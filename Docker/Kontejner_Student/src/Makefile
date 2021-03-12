# dnf install glibc-static glib2-static libstdc++-static upx

all: kontrolka ukladatko print_me upx
static: all

kontrolka: kontrolka.c
	cc -o kontrolka kontrolka.c -static
ukladatko: ukladatko.c
	cc -o ukladatko ukladatko.c -static
print_me: print_me.c
	cc -o print_me print_me.c -static
upx: kontrolka ukladatko print_me
	upx kontrolka
	upx ukladatko
	upx print_me

dynamic : kontrolka-dynamic ukladatko-dynamic print_me-dynamic

kontrolka-dynamic: kontrolka.c
	cc -o kontrolka kontrolka.c
ukladatko-dynamic: ukladatko.c
	cc -o ukladatko ukladatko.c
print_me-dynamic: print_me.c
	cc -o print_me print_me.c

clean:
	rm -f kontrolka ukladatko print_me