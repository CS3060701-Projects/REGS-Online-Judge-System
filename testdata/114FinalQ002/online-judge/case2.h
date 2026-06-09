#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "440e4e1e5ebf974baf461885f3eb6246b457afdbf0b37860f91f0d344a38d96d"

#include <sstream>
#include <string>

#include "command.h"
#include "test.h"

inline std::string trimEnd(std::string s) {
  s.erase(s.find_last_not_of("\t\n\r\f\v ") + 1);
  return s;
}

inline std::string trimEnd(std::ostringstream& oss) {
  auto s = oss.str();
  s.erase(s.find_last_not_of("\t\n\r\f\v ") + 1);
  return s;
}

TEST_CASE("case2") {
  std::istringstream inputCommand(
      "Insert \"Database Systems\" \"Ramez\" 6 2\n"
      "Insert \"Database Systems\" \"Ramez\" 7 3\n"
      "Insert \"Algorithms\" \"CLRS\" 3 1\n"
      "Insert \"Algorithms\" \"CLRS\" 4 2\n"
      "Insert \"Clean Code\" \"Robert Martin\" 1 1\n"
      "Insert \"Design Patterns\" \"GoF\" 1 2\n"
      "Insert \"Design Patterns\" \"GoF\" 2 1\n"
      "Insert \"Compilers\" \"Aho\" 2 3\n"
      "Insert \"Temp Book\" \"Author X\" 1 1\n"
      "Find Author \"CLRS\"\n"
      "Find Author \"Nobody\"\n"
      "Borrow \"Algorithms\" \"CLRS\" 3 \"Tom\"\n"
      "Borrow \"Algorithms\" \"CLRS\" 4 \"Tom\"\n"
      "Borrow \"Clean Code\" \"Robert Martin\" 1 \"Tom\"\n"
      "Borrow \"Clean Code\" \"Robert Martin\" 1 \"Jerry\"\n"
      "Borrow \"Design Patterns\" \"GoF\" 1 \"Tom\"\n"
      "Borrow \"Design Patterns\" \"GoF\" 1 \"Jerry\"\n"
      "Borrow \"Design Patterns\" \"GoF\" 2 \"Sam\"\n"
      "Delete Edition \"Design Patterns\" \"GoF\" 1\n"
      "Delete Edition \"Design Patterns\" \"GoF\" 2\n"
      "Delete Edition \"Algorithms\" \"CLRS\" 3\n"
      "Delete Edition \"Algorithms\" \"CLRS\" 4\n"
      "Delete Edition \"NoExist\" \"Nobody\" 1\n"
      "Return \"Algorithms\" \"CLRS\" 3 \"Tom\"\n"
      "Return \"Algorithms\" \"CLRS\" 4 \"Tom\"\n"
      "Delete Edition \"Algorithms\" \"CLRS\" 3\n"
      "Delete Edition \"Algorithms\" \"CLRS\" 4\n"
      "Insert \"Algorithms\" \"CLRS\" 3 5\n"
      "Find Book \"Algorithms\" \"CLRS\"\n"
      "Sort by Author\n"
      "Borrow \"Database Systems\" \"Ramez\" 6 \"Alice\"\n"
      "Borrow \"Database Systems\" \"Ramez\" 6 \"Bob\"\n"
      "Borrow \"Database Systems\" \"Ramez\" 6 \"Charlie\"\n"
      "Find Book \"Database Systems\" \"Ramez\"\n"
      "Return \"Database Systems\" \"Ramez\" 6 \"Alice\"\n"
      "Find Book \"Database Systems\" \"Ramez\"\n"
      "Borrow \"Compilers\" \"Aho\" 2 \"Alice\"\n"
      "Borrow \"Compilers\" \"Aho\" 2 \"Bob\"\n"
      "Borrow \"Compilers\" \"Aho\" 2 \"Charlie\"\n"
      "Borrow \"Compilers\" \"Aho\" 2 \"David\"\n"
      "Borrow \"Algorithms\" \"CLRS\" 3 \"Zoe\"\n"
      "Delete Edition \"Algorithms\" \"CLRS\" 3\n"
      "Return \"Algorithms\" \"CLRS\" 3 \"Zoe\"\n"
      "Delete Edition \"Algorithms\" \"CLRS\" 3\n"
      "Delete Edition \"Temp Book\" \"Author X\" 1\n"
      "Find Book \"Temp Book\" \"Author X\"\n"
      "Sort by Title\n"
      "Return \"Clean Code\" \"Robert Martin\" 1 \"Tom\"\n"
      "Delete Edition \"Clean Code\" \"Robert Martin\" 1\n");
  std::ostringstream oss;
  std::streambuf* cout_buf = std::cout.rdbuf();
  std::streambuf* cin_buf = std::cin.rdbuf();
  std::cout.rdbuf(oss.rdbuf());
  std::cin.rdbuf(inputCommand.rdbuf());
  CHECK(true);

  string line;
  while (getline(cin, line)) {
    execute(line);
  }

  auto expected =
      "Insert Ramez's Database Systems, Edition: 6.;"
      "Insert Ramez's Database Systems, Edition: 7.;"
      "Insert CLRS's Algorithms, Edition: 3.;"
      "Insert CLRS's Algorithms, Edition: 4.;"
      "Insert Robert Martin's Clean Code, Edition: 1.;"
      "Insert GoF's Design Patterns, Edition: 1.;"
      "Insert GoF's Design Patterns, Edition: 2.;"
      "Insert Aho's Compilers, Edition: 2.;"
      "Insert Author X's Temp Book, Edition: 1.;"
      "CLRS's Books: <Algorithms>;"
      "No book found.;"
      "Tom borrows CLRS's Algorithms, Edition: 3.;"
      "Tom borrows CLRS's Algorithms, Edition: 4.;"
      "Tom borrows Robert Martin's Clean Code, Edition: 1.;"
      "No available copies.;"
      "Tom borrows GoF's Design Patterns, Edition: 1.;"
      "Jerry borrows GoF's Design Patterns, Edition: 1.;"
      "Sam borrows GoF's Design Patterns, Edition: 2.;"
      "Outstanding borrows exist.;"
      "Outstanding borrows exist.;"
      "Outstanding borrows exist.;"
      "Outstanding borrows exist.;"
      "Datum doesn't exist.;"
      "Tom returns CLRS's Algorithms, Edition: 3.;"
      "Tom returns CLRS's Algorithms, Edition: 4.;"
      "Delete CLRS's Algorithms, Edition: 3.;"
      "Delete CLRS's Algorithms, Edition: 4.;"
      "Insert CLRS's Algorithms, Edition: 3.;"
      "Title: Algorithms	Author: CLRS	Edition: <3(5)>;"
      "Title: Compilers	Author: Aho	Edition: <2(3)>;"
      "Title: Temp Book	Author: Author X	Edition: <1(1)>;"
      "Title: Algorithms	Author: CLRS	Edition: <3(5)>;"
      "Title: Design Patterns	Author: GoF	Edition: <1(0), 2(0)>;"
      "Title: Database Systems	Author: Ramez	Edition: <6(2), 7(3)>;"
      "Title: Clean Code	Author: Robert Martin	Edition: <1(0)>;"
      "Alice borrows Ramez's Database Systems, Edition: 6.;"
      "Bob borrows Ramez's Database Systems, Edition: 6.;"
      "No available copies.;"
      "Title: Database Systems	Author: Ramez	Edition: <6(0), 7(3)>;"
      "Alice returns Ramez's Database Systems, Edition: 6.;"
      "Title: Database Systems	Author: Ramez	Edition: <6(1), 7(3)>;"
      "Alice borrows Aho's Compilers, Edition: 2.;"
      "Bob borrows Aho's Compilers, Edition: 2.;"
      "Charlie borrows Aho's Compilers, Edition: 2.;"
      "No available copies.;"
      "Zoe borrows CLRS's Algorithms, Edition: 3.;"
      "Outstanding borrows exist.;"
      "Zoe returns CLRS's Algorithms, Edition: 3.;"
      "Delete CLRS's Algorithms, Edition: 3.;"
      "Delete Author X's Temp Book, Edition: 1.;"
      "Book doesn't exist.;"
      "Title: Clean Code	Author: Robert Martin	Edition: <1(0)>;"
      "Title: Compilers	Author: Aho	Edition: <2(0)>;"
      "Title: Database Systems	Author: Ramez	Edition: <6(1), 7(3)>;"
      "Title: Design Patterns	Author: GoF	Edition: <1(0), 2(0)>;"
      "Tom returns Robert Martin's Clean Code, Edition: 1.;"
      "Delete Robert Martin's Clean Code, Edition: 1.;";

  std::cout.rdbuf(cout_buf);
  std::cin.rdbuf(cin_buf);
  CHECK(trimEnd(oss) == trimEnd(expected));
}

#endif
