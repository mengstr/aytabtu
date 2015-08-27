#include <fcntl.h>
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>



#define RTTSTATS_FILENAME       "rtt.dat"

typedef struct {
    uint8_t ms[320];
} rttstat_t;


//
//
//
void UpdateRttStats(uint16_t ms) {
    rttstat_t data;
    int f;
    int i;

    // Scale and cap the rtt so it fits in a uint8
    if (ms < 800) ms = 800;
    ms -= 800;
    ms /= 10;
    if (ms > 220) ms = 220;
    printf("Scaled value is %d\n",ms);

    // Try to open the file and if it fails create a new file
    // and popuate it with an empty record for usage further down
    f = open(RTTSTATS_FILENAME, O_RDWR);
    if (f < 1) {
        f = open(RTTSTATS_FILENAME, O_RDWR | O_CREAT);
        memset(&data, 255, sizeof(data));
    } else {
        read(f, (char *)&data, sizeof(data));
    }

    for (i = 319; i > 0; i--) {
        data.ms[i] = data.ms[i - 1];
    }
    data.ms[0] = ms;
    lseek(f, 0, SEEK_SET);
    write(f, (char *)&data, sizeof(data));
    close(f);

    return;
}


//
//
//
void DisplayRtt(int mode) {
    rttstat_t data;
    int f;
    int x;

    f = open(RTTSTATS_FILENAME, O_RDWR);
    if (f < 1) {
        return;
    }
    read(f, (char *)&data, sizeof(data));
    close(f);

    for (x = 0; x < 320; x++) {
        data.ms[x] = x / 2;
        if (data.ms[x] <= 220 ) {
            printf("[%d] = %d\n",x,data.ms[x]);
        }
    }

}




int main(int argc, char *argv[]) {
    uint16_t v;

    if (argc<2) return 1;
    v=atol(argv[1]);
    printf("Inserting %d into file\n",v);
    UpdateRttStats(v);
//    DisplayRtt(0);    
    return 0;
}

