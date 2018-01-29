cmake_minimum_required (VERSION 3.10)
project (h2o)

if(CMAKE_BUILD_TYPE STREQUAL "Release")    
    set(CMAKE_CXX_LINK_FLAGS "${CMAKE_CXX_LINK_FLAGS} -s")
endif()

AUX_SOURCE_DIRECTORY(src DIR_SRCS)
ADD_EXECUTABLE(h2o ${DIR_SRCS})
