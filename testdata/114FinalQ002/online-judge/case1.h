#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "4bf3f651fce61b407506c0e26b1279bbfb7fe8dfec2c493d9b5a7ef94cfb22d1"

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

TEST_CASE("case1") {
  std::istringstream inputCommand(
      "Insert \"The Great Gatsby\" \"F. Scott Fitzgerald\" 1 5\n"
      "Insert \"The Great Gatsby\" \"F. Scott Fitzgerald\" 1 5\n"
      "Insert \"1984\" \"George Orwell\" 1 3\n"
      "Insert \"1984\" \"George Orwell\" 2 2\n"
      "Insert \"Brave New World\" \"Aldous Huxley\" 1 2\n"
      "Insert \"Brave New World\" \"Aldous Huxley\" 3 1\n"
      "Insert \"Animal Farm\" \"George Orwell\" 1 4\n"
      "Find Book \"1984\" \"George Orwell\"\n"
      "Find Book \"Hamlet\" \"Shakespeare\"\n"
      "Find Author \"George Orwell\"\n"
      "Find Author \"Nobody\"\n"
      "Sort by Title\n"
      "Sort by Author\n"
      "Delete Edition \"1984\" \"George Orwell\" 99\n"
      "Delete Edition \"Hamlet\" \"Shakespeare\" 1\n"
      "Borrow \"1984\" \"George Orwell\" 1 \"Alice\"\n"
      "Borrow \"1984\" \"George Orwell\" 1 \"Bob\"\n"
      "Borrow \"1984\" \"George Orwell\" 1 \"Charlie\"\n"
      "Borrow \"1984\" \"George Orwell\" 1 \"Alice\"\n"
      "Borrow \"1984\" \"George Orwell\" 1 \"David\"\n"
      "Borrow \"1984\" \"George Orwell\" 5 \"Alice\"\n"
      "Borrow \"The Great Gatsby\" \"F. Scott Fitzgerald\" 1 \"Alice\"\n"
      "Return \"1984\" \"George Orwell\" 1 \"Eve\"\n"
      "Return \"The Great Gatsby\" \"F. Scott Fitzgerald\" 1 \"Bob\"\n"
      "Delete Edition \"1984\" \"George Orwell\" 1\n"
      "Return \"1984\" \"George Orwell\" 1 \"Alice\"\n"
      "Delete Edition \"1984\" \"George Orwell\" 1\n"
      "Return \"1984\" \"George Orwell\" 1 \"Bob\"\n"
      "Return \"1984\" \"George Orwell\" 1 \"Charlie\"\n"
      "Delete Edition \"1984\" \"George Orwell\" 1\n"
      "Find Book \"1984\" \"George Orwell\"\n"
      "Delete Edition \"Brave New World\" \"Aldous Huxley\" 1\n"
      "Delete Edition \"Brave New World\" \"Aldous Huxley\" 3\n"
      "Find Author \"Aldous Huxley\"\n"
      "Insert \"1984\" \"George Orwell\" 1 10\n"
      "Find Book \"1984\" \"George Orwell\"\n"
      "Sort by Title\n"
      "insert \"test\" \"test\" 1 1\n"
      "Insert\n"
      "Insert \"only title\"");
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
      "Insert F. Scott Fitzgerald's The Great Gatsby, Edition: 1.;"
      "Datum already exist.;"
      "Insert George Orwell's 1984, Edition: 1.;"
      "Insert George Orwell's 1984, Edition: 2.;"
      "Insert Aldous Huxley's Brave New World, Edition: 1.;"
      "Insert Aldous Huxley's Brave New World, Edition: 3.;"
      "Insert George Orwell's Animal Farm, Edition: 1.;"
      "Title: 1984	Author: George Orwell	Edition: <1(3), 2(2)>;"
      "Book doesn't exist.;"
      "George Orwell's Books: <1984, Animal Farm>;"
      "No book found.;"
      "Title: 1984	Author: George Orwell	Edition: <1(3), 2(2)>;"
      "Title: Animal Farm	Author: George Orwell	Edition: <1(4)>;"
      "Title: Brave New World	Author: Aldous Huxley	Edition: "
      "<1(2), 3(1)>;"
      "Title: The Great Gatsby	Author: F. Scott Fitzgerald	"
      "Edition: "
      "<1(5)>;"
      "Title: Brave New World	Author: Aldous Huxley	Edition: "
      "<1(2), 3(1)>;"
      "Title: The Great Gatsby	Author: F. Scott Fitzgerald	"
      "Edition: "
      "<1(5)>;"
      "Title: 1984	Author: George Orwell	Edition: <1(3), 2(2)>;"
      "Title: Animal Farm	Author: George Orwell	Edition: <1(4)>;"
      "Datum doesn't exist.;"
      "Datum doesn't exist.;"
      "Alice borrows George Orwell's 1984, Edition: 1.;"
      "Bob borrows George Orwell's 1984, Edition: 1.;"
      "Charlie borrows George Orwell's 1984, Edition: 1.;"
      "Already borrowed.;"
      "No available copies.;"
      "Datum doesn't exist.;"
      "Alice borrows F. Scott Fitzgerald's The Great Gatsby, Edition: 1.;"
      "Borrow record not found.;"
      "Borrow record not found.;"
      "Outstanding borrows exist.;"
      "Alice returns George Orwell's 1984, Edition: 1.;"
      "Outstanding borrows exist.;"
      "Bob returns George Orwell's 1984, Edition: 1.;"
      "Charlie returns George Orwell's 1984, Edition: 1.;"
      "Delete George Orwell's 1984, Edition: 1.;"
      "Title: 1984	Author: George Orwell	Edition: <2(2)>;"
      "Delete Aldous Huxley's Brave New World, Edition: 1.;"
      "Delete Aldous Huxley's Brave New World, Edition: 3.;"
      "No book found.;"
      "Insert George Orwell's 1984, Edition: 1.;"
      "Title: 1984	Author: George Orwell	Edition: <1(10), 2(2)>;"
      "Title: 1984	Author: George Orwell	Edition: <1(10), 2(2)>;"
      "Title: Animal Farm	Author: George Orwell	Edition: <1(4)>;"
      "Title: The Great Gatsby	Author: F. Scott Fitzgerald	"
      "Edition: "
      "<1(4)>;"
      "Unknown Command.;"
      "Incomplete Command.;"
      "Incomplete Command.;";

  std::cout.rdbuf(cout_buf);
  std::cin.rdbuf(cin_buf);
  CHECK(trimEnd(oss) == trimEnd(expected));
}

#endif
