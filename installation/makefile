CC=gcc 
SRC_DIR = ../bot_src
OBJ_DIR = ../bots

# Add you C bot name here 
_CSRCS = cbot under_500

CSRCS = $(patsubst %,$(SRC_DIR)/%,$(_CSRCS))

all: $(CSRCS)

$(CSRCS): %: %.c
	$(CC) -o $(subst $(SRC_DIR),$(OBJ_DIR),$@.exe) $< -lcjson	

.PHONY: clean 

clean: 
	rm -f $(OBJ_DIR)/*.exe 
