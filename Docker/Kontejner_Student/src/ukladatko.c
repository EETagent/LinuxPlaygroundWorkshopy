#include <stdio.h>
#include <dirent.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <string.h>
#define BUFSIZE 128

int parse_output(char *cmd1, const char *cmd2)
{
  char cmd[1024];

  snprintf(cmd, sizeof(cmd), "/bin/bash -c '(diff -q <(%s) <(%s)) &>/dev/null'", cmd1, cmd2);

  char buf[BUFSIZE];
  FILE *fp;

  if ((fp = popen(cmd, "r")) == NULL)
  {
    printf("Chyba při otevírání roury!\n");
    return -1;
  }

  while (fgets(buf, BUFSIZE, fp) != NULL)
  {
    if (strcmp(buf, "") != 0)
    {
      return -1;
    }
  }

  if (pclose(fp))
  {
    return -1;
  }

  return 0;
}

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

int main()
{
  double complete = 0;
  double all = 0;

  const char *files[7][256];
  files[0][0] = "cat ";
  files[0][1] = "/home/eidam/kyber/uloha1.txt";
  files[0][2] = "id -u";

  files[1][0] = "cat ";
  files[1][1] = "/home/eidam/kyber/uloha2.txt";
  files[1][2] = "/home/eidam/print_me";

  files[2][0] = "cat ";
  files[2][1] = "/home/eidam/kyber/uloha3.txt";
  files[2][2] = "hostname";

  files[3][0] = "cat ";
  files[3][1] = "/home/eidam/kyber/uloha4.txt";
  files[3][2] = "curl http://bezpecnost.ssps.cz/www/share/Hackdays/remote.txt";

  files[4][0] = "cat ";
  files[4][1] = "/home/eidam/kyber/uloha5.txt";
  files[4][2] = "echo 'vlajka{TLIWJXNCUJ}'";

  files[5][0] = "cat ";
  files[5][1] = "/home/eidam/kyber/uloha6.txt";
  files[5][2] = "echo 'BUDPRIPRAVEN'";

  files[6][0] = "cat ";
  files[6][1] = "/home/eidam/kyber/uloha7.txt";
  files[6][2] = "dig +short bezpecnost.ssps.cz";

  int filesLength = sizeof(files) / sizeof(files[0]);
  char filecmd[1048];

  for (int a = 0; a < filesLength; a = a + 1)
  {
    all = all + 1;

    int isOK = 0;
    strcpy(filecmd, files[a][0]);
    strcat(filecmd, files[a][1]);

    if (f_exists(files[a][1]) == 0)
    {
      if (parse_output(filecmd, files[a][2]) == 0)
      {
        isOK = 1;
      }
    }

    if (isOK == 1)
    {
      printf("OK -> \t");
      complete = complete + 1;
    }
    else
    {
      printf(":(( -> \t");
    }
    printf("%s\n", files[a][1]);
  }

  printf("\nDokončeno: %f %%\n", (complete / all) * 100);
  return 0;
}
