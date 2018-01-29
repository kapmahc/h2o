cmake_minimum_required (VERSION 3.10)
project (h2o)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

if(CMAKE_BUILD_TYPE STREQUAL "Release")
    set(CMAKE_CXX_LINK_FLAGS "${CMAKE_CXX_LINK_FLAGS} -s")
endif()

AUX_SOURCE_DIRECTORY(src DIR_SRCS)

ADD_EXECUTABLE(h2o ${DIR_SRCS})
TARGET_LINK_LIBRARIES(h2o wthttp wt boost_system boost_thread boost_filesystem boost_program_options pqxx)
