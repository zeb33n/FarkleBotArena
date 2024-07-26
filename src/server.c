#include <fcntl.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/poll.h>
#include <sys/socket.h>
#include <unistd.h>

struct Player {
  char num;
  int clientfd;
};

static char* read_line(int pipefd) {
  unsigned char c[] = {0};
  char* out = malloc(1024 * sizeof(char));
  int length = 1024;
  int i = 0;
  while (c[0] != ';') {
    if (c[0] == ';') {
      printf("hell yeah\n");
    }
    int bytes_read = read(pipefd, c, 1);
    if (bytes_read < 0) {
      perror("ReadError");
      exit(-1);
    }
    out[i] = *c;
    i++;
  }
  printf("%i\n", i);
  out[i] = '\0';
  return out;
}

int send_game_state(int pipefd, int clientfd) {
  char* game_state = read_line(pipefd);

  // int bytes_read = read(pipefd, game_state, 255);
  // if (bytes_read < 0) {
  //   perror("read");
  //   return -1;
  // }
  // game_state[bytes_read] = '\0';

  long send_err = send(clientfd, game_state, 255, 0);
  printf("%s\n", game_state);
  free(game_state);
  if (send_err == -1) {
    perror("Send Error");
    return -1;
  }

  return 0;
}

int recv_client_input(int clientfd) {
  char buffer[1] = {0};
  if (recv(clientfd, buffer, 1, 0) == 0) {
    printf("exiting\n");
    return 1;
  }
  if (*buffer == 'q') {
    printf("exiting\n");
    return 1;
  }
  int playerfd = open("../player0", O_WRONLY);
  write(playerfd, buffer, 1);
  printf("%s\n", buffer);
  close(playerfd);
  return 0;
}

void cp_file(char* source, char* target) {
  int srcfd = open(source, O_RDONLY);
  int tgtfd = open(target, O_CREAT | O_WRONLY, 0777);
  char buffer[4096];

  for (;;) {
    int n = read(srcfd, buffer, 4096);
    if (n < 0) {
      perror("Error reading file");
      exit(n);
    } else if (n == 0) {
      break;
    }

    int write_err = write(tgtfd, buffer, n);
    if (write_err == -1) {
      perror("Write Error:");
      exit(-1);
    }
  }
  close(srcfd);
  close(tgtfd);
}

void register_player(char num) {
  char target[] = "../bots/player0.py";
  target[14] = num;
  cp_file("../pysrc/player.py", target);
}

// TODO! use fork for multiple connections!
struct Player register_players(int socketfd) {
  int clientfd = accept(socketfd, 0, 0);
  if (clientfd == -1) {
    perror("accept error");
  }
  printf("connection recieved\n");
  register_player('0');
  struct Player player = {'0', clientfd};
  return player;
}

int await_game_start(struct Player player) {
  int sendfd = open("start", O_WRONLY);
  char* out = "1";
  struct pollfd fd[] = {{player.clientfd, POLLIN, 0}};
  while (1) {
    poll(fd, 1, 50000);
    if (fd[0].revents & POLLIN) {
      int err = write(sendfd, out, 1);
      close(sendfd);
      if (err == -1) {
        perror("writeerror");
        return -1;
      }
      return 0;
    }
  }
}

// pipe into file with player tag

// players join into lobby
// when lobby is populated send message to game loop to start
// ie when all players have readyed

int main() {
  int socketfd = socket(AF_INET, SOCK_STREAM, 0);
  struct sockaddr_in address = {AF_INET, htons(4123), INADDR_ANY};

  int bind_err = bind(socketfd, &address, sizeof(address));
  if (bind_err == -1) {
    perror("bind error");
    return 1;
  }

  int listen_err = listen(socketfd, 10);
  if (listen_err == -1) {
    perror("listen error");
    return 1;
  }

  struct Player player = register_players(socketfd);

  await_game_start(player);
  printf("starting\n");

  int pipefd = open("pipe", O_RDONLY);

  struct pollfd fds[2] = {{pipefd, POLLIN, 0}, {player.clientfd, POLLIN, 0}};
  for (;;) {
    poll(fds, 2, 50000);
    if (fds[0].revents & POLLIN) {
      int gs_err = send_game_state(pipefd, player.clientfd);
      if (gs_err == -1) {
        return -1;
      }

    } else if (fds[1].revents & POLLIN) {
      int exit = recv_client_input(player.clientfd);
      if (exit) {
        return exit;
      }
    }
  }

  close(pipefd);
  close(player.clientfd);
  close(socketfd);

  return 0;
}
