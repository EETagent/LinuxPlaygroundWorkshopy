#include <stddef.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <dirent.h>
#include <sys/types.h>
#include <sys/stat.h>
#define BUFSIZE 128

int f_exists(const char *fname)
{
	if (access(fname, F_OK) != -1)
	{
		return 0;
	}
	else
	{
		return 1;
	}
}

int dir_exists(const char *dirname)
{
	struct stat dirs;
	stat(dirname, &dirs);
	if (S_ISDIR(dirs.st_mode))
	{
		return 0;
	}
	else
	{
		return 1;
	}
}

int permissions(const char *fname, int value)
{
	struct stat perms;
	stat(fname, &perms);
	if (perms.st_mode == value)
	{
		return 0;
	}
	else
	{
		return 1;
	}
}

int main()
{
	double complete = 0;
	double all = 0;

	char *locations[11];
	locations[0] = "/home/eidam/adresar1";
	locations[1] = "/home/eidam/adresar1/adresar1_1";
	locations[2] = "/home/eidam/adresar1/adresar1_1/adresar1_1_1";
	locations[3] = "/home/eidam/adresar1/adresar1_1/adresar1_1_2";
	locations[4] = "/home/eidam/adresar1/adresar1_2";
	locations[5] = "/home/eidam/adresar1/adresar1_3";
	locations[6] = "/home/eidam/adresar1/adresar1_3/adresar1_3_1";
	locations[7] = "/home/eidam/adresar1/adresar1_3/adresar1_3_1/adresar1_3_1_1";
	locations[8] = "/home/eidam/adresar1/adresar1_3/adresar1_3_1/adresar1_3_1_2";
	locations[9] = "/home/eidam/adresar1/adresar1_4";
	locations[10] = "/home/eidam/adresar2";

	char *files[6];
	files[0] = "/home/eidam/adresar1/adresar1_1/soubor1.txt";
	files[1] = "/home/eidam/adresar1/adresar1_2/soubor2.txt";
	files[2] = "/home/eidam/adresar1/adresar1_2/soubor3.txt";
	files[3] = "/home/eidam/adresar1/adresar1_3/.soubor4.txt";
	files[4] = "/home/eidam/adresar1/adresar1_4/soubor5.txt";
	files[5] = "/home/eidam/adresar2/soubor6.txt";

	char *perms[3];
	perms[0] = "/home/eidam/adresar1/adresar1_2/soubor2.txt";
	perms[1] = "/home/eidam/adresar1/adresar1_2/soubor3.txt";
	perms[2] = "/home/eidam/adresar1/adresar1_3/.soubor4.txt";

	int values[3];
	values[0] = 33268; //100764 in octal
	values[1] = 33185; //100641 in octal
	values[2] = 33171; //100623 in octal

	//checkovani adresaru
	printf("Adrasáře:\n");
	int loc_length = sizeof(locations) / sizeof(locations[0]);
	for (int a = 0; a < loc_length; a++)
	{
		all = all + 1;
		if (dir_exists(locations[a]) == 0)
		{
			printf("OK ->\t");
			complete = complete + 1;
		}
		else
		{
			printf(":(( ->\t");
		}
		printf("%s\n", locations[a]);
	}

	//checkovani souboru
	printf("\nSoubory:\n");
	int f_length = sizeof(files) / sizeof(files[0]);
	for (int a = 0; a < f_length; a++)
	{
		all = all + 1;
		if (f_exists(files[a]) == 0)
		{
			printf("OK ->\t");
			complete = complete + 1;
		}
		else
		{
			printf(":(( ->\t");
		}
		printf("%s\n", files[a]);
	}

	//checkovani opravneni
	printf("\nOprávnění:\n");
	int perm_length = sizeof(perms) / sizeof(perms[0]);
	for (int a = 0; a < perm_length; a++)
	{
		all = all + 1;
		if (permissions(perms[a], values[a]) == 0)
		{
			printf("OK ->\t");
			complete = complete + 1;
		}
		else
		{
			printf(":(( ->\t");
		}
		printf("%s\n", perms[a]);
	}

	printf("\nDokončeno: %f %%\n", (complete / all) * 100);
	return 0;
}
