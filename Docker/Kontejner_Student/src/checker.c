#include <stdio.h>
#include <dirent.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <string.h>
#define BUFSIZE 128

int parse_output(char *cmd1, char *cmd2) {
    char cmd[1024];

    //printf("%s %s\n", cmd1, cmd2);
    snprintf(cmd, sizeof(cmd), "/bin/bash -c '(diff -q <(%s) <(%s)) &>/dev/null'", cmd1, cmd2);

    //printf("%s\n", cmd);

    char buf[BUFSIZE];
    FILE *fp;

    if ((fp = popen(cmd, "r")) == NULL) {
        printf("Error opening pipe!\n");
        return -1;
    }

    while (fgets(buf, BUFSIZE, fp) != NULL) {
        // Do whatever you want here...
        //printf("OUTPUT: %s", buf);
        if (buf != "") {
          return -1;
        }
    }

    if(pclose(fp))  {
        //rintf("Command not found or exited with error status\n");

        return -1;
    }

    return 0;
}

int exists(const char *fname)
{
  FILE *file;
  if( access( fname, F_OK ) != -1 ) {
    return 0;
  } else {
    return 1;
  }
}

int isDirectory(const char *path) {
   struct stat statbuf;
   if (stat(path, &statbuf) != 0) {
     return 0;
   }
   return 1;
}


int main()
{


  char home[64] = "/home/eidam/";

  double complete = 0;
  double all = 0;

  const char *locations[11];
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

  const char *files[8][256];
  files[0][0] = "cat ";
  files[0][1] = "/home/eidam/adresar1/adresar1_1/soubor1.txt";
  files[0][2] = "id -u";

  files[1][0] = "cat ";
  files[1][1] = "/home/eidam/adresar1/adresar1_2/soubor2.txt";
  files[1][2] = "/home/eidam/print_me";

  files[2][0] = "cat ";
  files[2][1] = "/home/eidam/adresar1/adresar1_2/soubor3.txt";
  files[2][2] = "hostname";

  files[3][0] = "cat ";
  files[3][1] = "/home/eidam/adresar1/adresar1_3/soubor4.txt";
  files[3][2] = "curl http://bezpecnost.ssps.cz/www/share/Hackdays/remote.txt";

  files[4][0] = "cat ";
  files[4][1] = "/home/eidam/adresar1/adresar1_4/soubor5.txt";
  files[4][2] = "echo 'vlajka{TLIWJXNCUJ}'";

  files[5][0] = "cat ";
  files[5][1] = "/home/eidam/adresar2/soubor6.txt";
  files[5][2] = "echo 'BUDPRIPRAVEN'";

  files[6][0] = "cat ";
  files[6][1] = "/home/eidam/adresar2/soubor7.txt";
  files[6][2] = "dig +short bezpecnost.ssps.cz";

  files[7][0] = "stat -c '%a' ";
  files[7][1] = "/home/eidam/adresar1/adresar1_3/adresar1_3_1/adresar1_3_1_1/";
  files[7][2] = "echo '751'";


  int locationLength = sizeof(locations) / sizeof(locations[0]);

  for(int a = 0; a < locationLength; a = a + 1 ){
    all = all + 1;

    int isOK = 0;
    if (isDirectory(locations[a]) != 0) {
      if (exists(locations[a]) == 0) {
        isOK = 1;
      }
    }

    if (isOK == 1) {
      printf("OK -> \t");
      complete = complete + 1;
    }
    else {
      printf("WRONG -> \t");
    }
    printf("%s\n", locations[a]);
   }


  int filesLength = sizeof(files) / sizeof(files[0]);
  char filecmd[1048];

  for (int a = 0; a < filesLength; a = a + 1) {
    all = all + 1;

    int isOK = 0;
    strcpy(filecmd, files[a][0]);
    strcat(filecmd, files[a][1]);

    //printf("%s\n",filecmd );
    //printf("%s\n",files[a][2]);

    if (exists(files[a][1]) == 0) {
      if (parse_output(filecmd, files[a][2]) == 0) {
        isOK = 1;
      }
    }

    if (isOK == 1) {
      printf("OK -> \t");
      complete = complete + 1;
    }
    else {
      printf("WRONG -> \t");
    }
    printf("%s\n", files[a][1]);
  }

  //printf("%i\n", (int)complete);
  //printf("%i\n", (int)all);

  printf("Done: %f %%\n", (complete / all) * 100);
  return 0;
}
