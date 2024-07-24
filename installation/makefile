CC=gcc 
BOT_SRC = ../bot_src
BOT_DIR = ../bots
SRC = ../src
TARGET = ../target

_CSRCS = server client
CSRCS = $(patsubst %,$(SRC)/%,$(_CSRCS))

# Add you C bot name here 
_BSRCS = cbot under_500
BSRCS = $(patsubst %,$(BOT_SRC)/%,$(_BSRCS))

all: sc bots

sc: $(CSRCS)

$(CSRCS): %: %.c
	$(CC) -o $(subst $(SRC),$(TARGET),$@) $< -Wno-incompatible-pointer-types

bots: $(BSRCS)

$(BSRCS): %: %.c
	$(CC) -o $(subst $(BOT_SRC),$(BOT_DIR),$@.exe) $< -lcjson	

.PHONY: clean 

clean: 
	rm -f $(BOT_DIR)/*.exe 
