#include <fcntl.h>
#include <netinet/in.h>
#include <stdio.h>
#include <sys/poll.h>
#include <sys/socket.h>
#include <unistd.h>

// char *read_pipe(int fd) {
//   static char out[256] = {0};
//   int bytes_read = read(fd, out, 255);
//   if (bytes_read < 0) {
//     perror("read");
//     return NULL;
//   }
//   out[bytes_read] = '\0';
//   return out;
// }

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
    return 1;
  }

  return 0;
}

int recv_client_input(int clientfd) {
  char buffer[256] = {0};
  if (recv(clientfd, buffer, 1, 0) == 0) {
    printf("exiting\n");
    return -1;
  }
  printf("%s\n", buffer);
  return 0;  
}

int register_player();

int pipe_player();

int main() {
  int pipefd = open("../pipe", O_RDONLY);
  int socketfd = socket(AF_INET, SOCK_STREAM, 0);
  struct sockaddr address = {AF_INET, htons(9999), 0};

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

  struct pollfd fds[2] = {{pipefd, POLLIN, 0}, {clientfd, POLLIN, 0}};
  for (;;) {
    poll(fds, 2, 50000);
    if (fds[0].revents & POLLIN) {
      int gs_err = send_game_state(pipefd, clientfd);
      if (gs_err == -1) {
        return -1;
      }

    } else if (fds[1].revents & POLLIN) {
      int exit = recv_client_input(clientfd); 
      if (exit) {
        return 0; 
      }
    }
  }

  close(clientfd);
  close(socketfd);

  return 0;
}
