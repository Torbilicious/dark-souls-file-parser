//
// Created by TBinder on 2019-01-18.
//

#ifndef DARK_SOULS_FILE_PARSER_READER_H
#define DARK_SOULS_FILE_PARSER_READER_H


#include <cstdio>

class Reader {

public:
    explicit Reader(const char *fileName);
    ~Reader();

    struct SlotDataHeaderStructure {
        unsigned int block_stat_size;
        unsigned int block_data_size;
    };

    struct SlotHeaderStructure {
        unsigned int block_metadata_high;
        unsigned int block_metadata_low;
        unsigned long block_size;
        unsigned int block_start_offset;
        unsigned int block_unknown_data_1;
        unsigned int block_skip_bytes;
        unsigned int end_of_block;
//        SlotDataHeaderStructure slot_data;
    };

    void printSaveFileStats();
    int getAmoundOfSlots();
    int getRealAmountOfSlots();

private:
    FILE *file;
    const long BLOCK_SIZE = 0x60190;
    const long BLOCK_INDEX = 0x2c0;
    const long BLOCK_DATA_OFFSET = 0x14;
    const long SLOTS_AMOUNT_OFFSET = 0xC;
    const long SLOTS_METADATA_OFFSET = 0x40;

    const long TIME_INDEX = 0x3c12c0;
    const long TIME_BLOCK_SIZE = 0x170;



};


#endif //DARK_SOULS_FILE_PARSER_READER_H
