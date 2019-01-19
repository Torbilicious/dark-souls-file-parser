//
// Created by TBinder on 2019-01-18.
//

#include <cstdlib>
#include <vector>
#include <iostream>
#include "Reader.h"

template <typename T>
T swap_endian(T u)
{
    static_assert (CHAR_BIT == 8, "CHAR_BIT != 8");

    union
    {
        T u;
        unsigned char u8[sizeof(T)];
    } source, dest;

    source.u = u;

    for (size_t k = 0; k < sizeof(T); k++)
        dest.u8[k] = source.u8[sizeof(T) - k - 1];

    return dest.u;
}

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

void Reader::doStuff() {
    for (int slotIndex = 0; slotIndex < getRealAmountOfSlots(); ++slotIndex) {


//        auto slotIndex = 0;
        auto offset = BLOCK_INDEX + BLOCK_SIZE * slotIndex;
        auto timeOffset = TIME_INDEX + TIME_BLOCK_SIZE * slotIndex;

        fseek(file, offset + 0x1f128, SEEK_SET);
        signed int deaths;
        std::fread(&deaths, sizeof deaths, 1, file);

        //    {'offset': 0x100, 'type': 'c', 'field': 'name', 'size': 14 * 2},


//        fseek(file, offset + 0x100, SEEK_SET);
//        std::vector<wchar_t> sName(14 * 2); // char is trivally copyable
//        std::fread(&sName, sizeof sName[0], 14 * 2, file);





        short stringlength = 14*2;
        wchar_t sName[stringlength]; //Or you can use malloc() / new instead.
        fseek(file , offset +0x100, SEEK_SET);
        std::fread(&sName[0], sizeof(wchar_t), 1, file);


        printf("sName:   %s", swap_endian(sName));
        printf("\n");
        printf("deaths: %d\n", deaths);
    }

}
