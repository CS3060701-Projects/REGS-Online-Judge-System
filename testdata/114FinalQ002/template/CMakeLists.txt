cmake_minimum_required(VERSION 3.20)

project(114FinalQ002)
enable_testing()
set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

define_property(GLOBAL
    PROPERTY PROBLEM_SECRETS
    BRIEF_DOCS "The secret define in test-case."
    FULL_DOCS "This value dynamic generate, use `configure_file` inject to case*.h, to avoid student skip unit test."
)

define_property(GLOBAL
    PROPERTY TEST_COUNTER
    BRIEF_DOCS "Test Configure count, just output case"
    FULL_DOCS "Test Configure count, just output case"
)
set_property(GLOBAL PROPERTY TEST_COUNTER 1)

get_filename_component(PROBLEM_ROOT_ABS "${CMAKE_CURRENT_SOURCE_DIR}/.." ABSOLUTE)
set(PROBLEM_ROOT "${PROBLEM_ROOT_ABS}" CACHE PATH "Path to the problem root directory")
set(SOURCE_DIR "${CMAKE_CURRENT_SOURCE_DIR}" CACHE PATH "Path to the current source directory")
set(SPEC_DIR "${PROBLEM_ROOT}/spec" CACHE PATH "Path to the spec directory")


message("Collect source files")
file(GLOB CXX_SOURCE_FILES
  "${SOURCE_DIR}/*.cxx"
  "${SOURCE_DIR}/*.cpp"
  "${SOURCE_DIR}/*.cc"
  "${SOURCE_DIR}/*.c"
)
foreach(SOURCE_FILE ${CXX_SOURCE_FILES})
  message(STATUS "Add : ${SOURCE_FILE}")
endforeach()

function(AddJudge case_file)
  get_property(count GLOBAL PROPERTY TEST_COUNTER)
  get_property(secrets GLOBAL PROPERTY PROBLEM_SECRETS)
  math(EXPR next_count "${count} + 1")
  
  set(TARGET_NAME "${case_file}")
  set(BUILDDIR "${CMAKE_CURRENT_BINARY_DIR}/case/_${count}_")
  
  set(SPEC_SRC "${SPEC_DIR}/${case_file}.h")
  set(SPEC_DST "${BUILDDIR}/case.h")
  set(SPEC_VAR "${PROBLEM_ROOT}/online-judge/${case_file}.h")

  if(NOT EXISTS "${SPEC_SRC}")
    message(WARNING "${SPEC_SRC} does not exist, skip")
    return()
  endif()

  set(SECRET "@SECRET@")
  file(MAKE_DIRECTORY "${BUILDDIR}")
  message("Create ${TARGET_NAME} Build Context")
  configure_file("${SPEC_SRC}" "${SPEC_DST}" @ONLY)
  message(STATUS "Copy ${SPEC_SRC}")
  message(STATUS "To ${SPEC_DST}")

  add_executable(${TARGET_NAME} ${CXX_SOURCE_FILES})
  target_include_directories(${TARGET_NAME} BEFORE PRIVATE
    "${BUILDDIR}"
  )

  message(STATUS "Add include: ${SOURCE_DIR}")
  message(STATUS "Add include: ${SOURCE_DIR}/include")
  target_include_directories(${TARGET_NAME} PRIVATE
    "${SOURCE_DIR}"
    "${SOURCE_DIR}/include"
  )

  message(STATUS "Generated target ${TARGET_NAME} using ${case_file}.h")
  add_test(NAME "${TARGET_NAME}" COMMAND ${TARGET_NAME})
  message(STATUS "Added test target: ${TARGET_NAME}")

  list(APPEND secrets "${case_file} '${SECRET}'")
  set_property(GLOBAL PROPERTY TEST_COUNTER ${next_count})
  set_property(GLOBAL PROPERTY PROBLEM_SECRETS ${secrets})
endfunction()

AddJudge("case1")
AddJudge("case2")
AddJudge("case3")
AddJudge("case4")
AddJudge("case5")

get_property(secrets GLOBAL PROPERTY PROBLEM_SECRETS)
message("All Secret:")
foreach(secret ${secrets})
  message(STATUS " ${secret}")
endforeach()
