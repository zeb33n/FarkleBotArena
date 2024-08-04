#include <fcntl.h>
#include <netinet/in.h>
#include <pthread.h>  //maybe use threads instead
#include <stdio.h>
#include <stdlib.h>
#include <sys/poll.h>
#include <sys/socket.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

struct Player {
  int clientfd;
  char num;
  char spipename[48];
  char rpipename[48];
};

int send_game_state(int pipefd, int clientfd) {
  printf("noice\n");
  char game_state[1024] = {0};

  char c = 0;
  int i = 0;
  while (c != ';' && read(pipefd, &c, 1) > 0) {
    game_state[i] = c;
    i++;
  }
  game_state[i - 1] = '\0';

  int send_err = send(clientfd, game_state, 1023, 0);
  if (send_err == -1) {
    perror("Send Error");
    exit(send_err);
  }

  return 0;
}

int recv_client_input(char* spipename, int clientfd) {
  char buffer[1] = {0};
  if (recv(clientfd, buffer, 1, 0) == 0) {
    printf("exiting\n");
    return 1;
  }
  if (*buffer == 'q') {
    printf("exiting\n");
    return 1;
  }
  int playerfd = open(spipename, O_WRONLY);
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

struct Player register_player(int clientfd, char id) {
  char playername[11] = {0};
  char prpipe[48] = {0};
  char pspipe[48] = {0};
  char target[48] = {0};
  sprintf(playername, "player%i", id);
  sprintf(prpipe, "../pipes/r%s", playername);
  sprintf(pspipe, "../pipes/s%s", playername);
  sprintf(target, "../bots/%s.py", playername);
  int err = mkfifo(prpipe, 0666);
  if (err < 0) {
    perror("fifo error");
    exit(err);
  }
  err = mkfifo(pspipe, 0666);
  if (err < 0) {
    perror("fifo error");
    exit(err);
  }
  cp_file("../pysrc/player.py", target);
  struct Player player = {clientfd, id, *pspipe, *prpipe};
  return player;
}

// TODO! use fork for multiple connections!
struct Player register_players(int socketfd) {
  char number_of_conn = 0;
  int clientfd;
  for (;;) {
    clientfd = accept(socketfd, 0, 0);
    if (clientfd == -1) {
      perror("accept error");
      continue;
    }
    int pid = fork();
    number_of_conn++;
    if (pid != 0) {
      break;
    }
  }
  struct Player player = register_player(clientfd, number_of_conn);
  printf("connection recieved\n");
  return player;
}

int await_game_start(struct Player player) {
  if (player.num != 1) {
    return 0;
  }
  // //int err = send(player.clientfd, "hello player 1", 48, 0);
  // if (err == -1) {
  //   perror("Send Error:");
  //   exit(err);
  // }
  int sendfd = open("start", O_WRONLY);
  char* out = "1";
  struct pollfd fd[] = {{player.clientfd, POLLIN, 0}};
  printf("3\n");
  while (1) {
    int err = poll(fd, 1, 50000);
    if (err == -1) {
      perror("Poll error");
    }
    if (fd[0].revents & POLLIN) {
      printf("bananas\n");
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

  printf("playerid: %i registered\n", player.num);

  await_game_start(player);

  int rpipefd = open(player.rpipename, O_RDONLY);

  struct pollfd fds[2] = {{rpipefd, POLLIN, 0}, {player.clientfd, POLLIN, 0}};
  for (;;) {
    int err = poll(fds, 2, 50000);
    if (err == -1) {
      perror("Poll error");
      return err;
    }
    if (fds[0].revents & POLLIN) {
      printf("here\n");
      int gs_err = send_game_state(rpipefd, player.clientfd);
      if (gs_err == -1) {
        return -1;
      }

    } else if (fds[1].revents & POLLIN) {
      printf("there\n");
      int exit = recv_client_input(player.spipename, player.clientfd);
      if (exit) {
        return exit;
      }
    }
  }

  close(rpipefd);
  close(player.clientfd);
  close(socketfd);

  return 0;
}
