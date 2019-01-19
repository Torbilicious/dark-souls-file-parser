#include <iostream>
#include <vector>
#include "Reader.h"


int main() {
//    const char *fileName = "/Volumes/Secured/Projects/cpp/dark-souls-file-parser/resources/DRAKS0005.sl2";
    const char *fileName = "/Volumes/Secured/Projects/cpp/dark-souls-file-parser/resources/DRAKS0005.sl2.2";

    auto reader = new Reader(fileName);
    reader->printSaveFileStats();

    return 0;
}