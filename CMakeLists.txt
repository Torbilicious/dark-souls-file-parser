cmake_minimum_required(VERSION 3.13)
project(dark_souls_file_parser)

set(CMAKE_CXX_STANDARD 17)

add_executable(dark_souls_file_parser main.cpp Reader.cpp Reader.h)