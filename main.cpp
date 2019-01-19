#include <iostream>
#include <vector>
#include "Reader.h"


int main() {
//    const char *fileName = "/Volumes/Secured/Projects/cpp/dark-souls-file-parser/DRAKS0005.sl2";
    const char *fileName = "/Volumes/Secured/Projects/cpp/dark-souls-file-parser/DRAKS0005.sl2.2";

    auto reader = new Reader(fileName);


//    auto slotAmount = reader->getAmoundOfSlots();
//    printf("Amount: %d\n", slotAmount);
//
//
//    auto slots = reader->getRealAmountOfSlots(slotAmount);
//
//    for (auto slot:slots) {
//        printf("block_metadata_high: %d\n", slot.block_metadata_high);
//    }

    reader->doStuff();

    return 0;
}