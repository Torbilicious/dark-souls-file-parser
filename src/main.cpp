#include <chrono>
#include "Reader.h"

using namespace std;
using namespace std::chrono;


int main() {
    const char *fileName = "/Volumes/Secured/Projects/cpp/dark-souls-file-parser/resources/DRAKS0005.sl2";
//    const char *fileName = "/Volumes/Secured/Projects/cpp/dark-souls-file-parser/resources/DRAKS0005.sl2.2";
//    const char *fileName = "/Volumes/Secured/Projects/cpp/dark-souls-file-parser/resources/DRAKS0005.sl2.3";

    auto reader = new Reader(fileName);


    auto running = true;
    auto counter = 0;
    auto maxTurns = 1;
    while (running) {
        reader->printSaveFileStats();

        running = ++counter < maxTurns;
    }

    return 0;
}