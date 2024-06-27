CC=gcc 
SRC_DIR = ../bot_src
OBJ_DIR = ../bots
_SRCS = cbot under_500
SRCS = $(patsubst %,$(SRC_DIR)/%,$(_SRCS))

all: $(SRCS)

$(SRCS): %: %.c
	$(CC) -o $(subst $(SRC_DIR),$(OBJ_DIR),$@.exe) $< -lcjson	

.PHONY: clean 

clean: 
	rm -f $(OBJ_DIR)/*.exe 
