//
// Created by TBinder on 2019-01-18.
//

#include <cstdlib>
#include <vector>
#include <iostream>
#include <wchar.h>
#include "Reader.h"


Reader::Reader(const char *fileName) {
    this->file = std::fopen(fileName, "r");

    if (file == nullptr) {
        printf("File could not be read!");
        exit(1);
    }
}

Reader::~Reader() {
    fclose(file);
}

int Reader::getAmoundOfSlots() {
    signed int slotAmount;
    fseek(file, SLOTS_AMOUNT_OFFSET, SEEK_SET);
    std::fread(&slotAmount, sizeof slotAmount, 1, file);

    return slotAmount;
}


int Reader::getRealAmountOfSlots() {
    int amountOfSlots = getAmoundOfSlots();

    SlotHeaderStructure slots[amountOfSlots];
    fseek(file, SLOTS_METADATA_OFFSET, SEEK_SET);
    std::fread(&slots, sizeof slots[0], amountOfSlots, file);

    std::vector<SlotHeaderStructure> v(slots, slots + sizeof slots / sizeof slots[0]);

    auto realSlots = 0;
    for (auto slot:v) {
        fseek(file, slot.block_start_offset + BLOCK_DATA_OFFSET, SEEK_SET);
        std::byte stuff;
        std::fread(&stuff, sizeof stuff, 1, file);

        std::byte badByte = {};
        if (stuff == badByte) continue;
        realSlots++;
    }

    return realSlots;
}

void Reader::printSaveFileStats() {
    for (int slotIndex = 0; slotIndex < getRealAmountOfSlots(); ++slotIndex) {
        auto offset = BLOCK_INDEX + BLOCK_SIZE * slotIndex;
        auto timeOffset = TIME_INDEX + TIME_BLOCK_SIZE * slotIndex;

        fseek(file, offset + 0x1f128, SEEK_SET);
        signed int deaths;
        std::fread(&deaths, sizeof deaths, 1, file);


        //    {'offset': 0x100, 'type': 'c', 'field': 'name', 'size': 14 * 2},
        //FIXME
        short stringlength = 14 * 2;
        wchar_t name[stringlength];
        fseek(file, offset + 0x100, SEEK_SET);
        std::fread(&name[0], sizeof(wchar_t), 1, file);

        wprintf(L"name:   %ls", name);
        printf("\n");
        printf("deaths: %d\n", deaths);
    }

}
