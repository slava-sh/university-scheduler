cmake_minimum_required(VERSION 3.7)
project(scheduler)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -Wextra")

if(INPUT_FILE)
  add_definitions(-DINPUT_FILE="${INPUT_FILE}")
endif()

include_directories(${PROJECT_SOURCE_DIR})

file(GLOB SOURCES "scheduler/*.cpp")

add_executable(scheduler ${SOURCES})
