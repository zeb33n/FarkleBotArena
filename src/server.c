#include <fcntl.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/poll.h>
#include <sys/socket.h>
#include <sys/stat.h>
#include <unistd.h>

int send_game_state(int pipefd, int clientfd) {
  char game_state[256] = {0};
  int bytes_read = read(pipefd, game_state, 255);
  if (bytes_read < 0) {
    perror("read");
    return -1;
  }
  game_state[bytes_read] = '\0';

  long send_err = send(clientfd, game_state, 255, 0);
  if (send_err == -1) {
    perror("Send Error");
    return -1;
  }

  return 0;
}

int recv_client_input(int clientfd, int outfd) {
  char buffer[1] = {0};
  if (recv(clientfd, buffer, 1, 0) == 0) {
    printf("exiting\n");
    return 1;
  }
  int write_err = write(outfd, buffer, 1);
  if (write_err != 1) {
    perror("Write Error");
    return -1;
  }
  close(outfd);
  return 0;
}

void cp_file(char *source, char *target) {
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

char *register_player(char num) {
  printf("1");
  char target[] = "../bots/player0.py";
  printf("2");
  static char pipe[] = "../player0";
  printf("3");
  target[14] = num;
  printf("4");
  pipe[9] = num;
  printf("5");
  int fifoerr = mkfifo(pipe, 0666);
  printf("6");
  printf("%i", fifoerr);
  if (fifoerr < 0) {
    perror("Fifo Error");
    exit(-1);
  }
  printf("7");
  cp_file("../player.py", target);
  printf("wow!");
  return pipe;
}

// cp script into dir
// make fifo for script

// player script when turned on needs to monitor for latest input.
// maybe add a turn counter so it know exactly what one.
// could be important with latency.

// block players from sending when its not there turn

int main() {

  int pipefd = open("../pipe", O_RDONLY);
  printf("wemadeit!"); 
  int socketfd = socket(AF_INET, SOCK_STREAM, 0);
  struct sockaddr address = {AF_INET, htons(9998), 0};

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

  int clientfd = accept(socketfd, 0, 0);
  if (clientfd == -1) {
    perror("accept error");
    return 1;
  }

  printf("connection recieved\n");
  char *player_pipe = register_player('0');
  int outfd = open(player_pipe, O_WRONLY);

  struct pollfd fds[2] = {{pipefd, POLLIN, 0}, {clientfd, POLLIN, 0}};
  for (;;) {
    poll(fds, 2, 50000);
    if (fds[0].revents & POLLIN) {
      printf("sending");
      int gs_err = send_game_state(pipefd, clientfd);
      if (gs_err == -1) {
        return -1;
      }

    } else if (fds[1].revents & POLLIN) {
      printf("recieving");
      int exit = recv_client_input(clientfd, outfd);
      if (exit) {
        return exit;
      }
    }
  }

  close(clientfd);
  close(socketfd);

  return 0;
}
