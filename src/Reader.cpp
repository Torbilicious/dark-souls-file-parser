//
// Created by TBinder on 2019-01-18.
//

#include <cstdlib>
#include <vector>
#include <iostream>
#include <wchar.h>
#include <string>
#include "Reader.h"


Reader::Reader(const char *fileName) {
    this->file = std::fopen(fileName, "r");

    if (file == nullptr) {
        printf("File could not be read!\n");
        exit(1);
    }

    if (!isCorrectSaveFile()) {
        printf("File is not in correct format.\n");
        exit(2);
    }
}

Reader::~Reader() {
    fclose(file);
}

bool Reader::isCorrectSaveFile() {
    const short stringlengthFmt = 4;
    std::string fmt( stringlengthFmt, '\0' );
    fseek(file , 0, SEEK_SET);
    fread(&fmt[0], sizeof(char), (size_t)stringlengthFmt, file);


    const short stringlengthVersion = 8;
    std::string version( stringlengthVersion, '\0' );
    fseek(file , 0x18, SEEK_SET);
    fread(&version[0], sizeof(char), (size_t)stringlengthVersion, file);


    printf("File header check:\n");
    printf("fmt:     %s\n", fmt.c_str());
    printf("version: %s\n", version.c_str());
    printf("\n");



    return fmt == "BND4" && version == "00000001";
}

int Reader::getAmoundOfSlots() {
    signed int slotAmount;
    fseek(file, SLOTS_AMOUNT_OFFSET, SEEK_SET);
    std::fread(&slotAmount, sizeof slotAmount, 1, file);

    return slotAmount;
}


int Reader::getRealAmountOfSlots() {
    const int amountOfSlots = getAmoundOfSlots();

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

        fseek(file, offset + 0x1f128, SEEK_SET);
        signed int deaths;
        std::fread(&deaths, sizeof deaths, 1, file);

        typedef wchar_t nameCharType;

        //    {'offset': 0x100, 'type': 'c', 'field': 'name', 'size': 14 * 2},
        //FIXME
        const short stringlength = 14 * 2;
        nameCharType name[stringlength];
        fseek(file, offset + 0x100, SEEK_SET);
        std::fread(&name[0], sizeof(nameCharType), 1, file);

        wprintf(L"name:   %ls", name);
        printf("\n");
        printf("deaths: %d\n", deaths);
    }

}
