#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "@SECRET@"

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

TEST_CASE("case5") {
  std::istringstream inputCommand(
      "Insert \"Alpha\" \"Zed\" 10 1\n"
      "Insert \"Alpha\" \"Ann\" 2 2\n"
      "Insert \"Alpha\" \"Ann\" 1 1\n"
      "Insert \"Beta\" \"Ann\" 3 1\n"
      "Insert \"Gamma Ray\" \"Ann\" 1 3\n"
      "Insert \"Delta\" \"Bob Smith\" 2020 2\n"
      "Insert \"Delta\" \"Bob Smith\" 2019 1\n"
      "Insert \"Omega\" \"Zed\" 0 1\n"
      "Insert \"Alpha\" \"Ann\" 2 99\n"
      "Find Book \"Alpha\" \"Ann\"\n"
      "Find Author \"Ann\"\n"
      "Sort by Title\n"
      "Sort by Author\n"
      "Borrow \"Alpha\" \"Ann\" 1 \"Carol\"\n"
      "Borrow \"Alpha\" \"Ann\" 1 \"Dave\"\n"
      "Borrow \"Alpha\" \"Ann\" 1 \"Carol\"\n"
      "Delete Edition \"Alpha\" \"Ann\" 1\n"
      "Borrow \"Alpha\" \"Ann\" 2 \"Carol\"\n"
      "Borrow \"Alpha\" \"Ann\" 2 \"Dave\"\n"
      "Borrow \"Alpha\" \"Ann\" 2 \"Eve\"\n"
      "Find Book \"Alpha\" \"Ann\"\n"
      "Return \"Alpha\" \"Ann\" 1 \"Carol\"\n"
      "Delete Edition \"Alpha\" \"Ann\" 1\n"
      "Find Book \"Alpha\" \"Ann\"\n"
      "Borrow \"Alpha\" \"Ann\" 1 \"Frank\"\n"
      "Return \"Alpha\" \"Ann\" 2 \"Eve\"\n"
      "Return \"Alpha\" \"Ann\" 2 \"Carol\"\n"
      "Delete Edition \"Alpha\" \"Ann\" 2\n"
      "Return \"Alpha\" \"Ann\" 2 \"Dave\"\n"
      "Delete Edition \"Alpha\" \"Ann\" 2\n"
      "Find Book \"Alpha\" \"Ann\"\n"
      "Find Author \"Ann\"\n"
      "Borrow \"Delta\" \"Bob Smith\" 2019 \"Carol\"\n"
      "Borrow \"Delta\" \"Bob Smith\" 2020 \"Carol\"\n"
      "Find Book \"Delta\" \"Bob Smith\"\n"
      "Sort by Title\n"
      "Return \"Delta\" \"Bob Smith\" 2019 \"Carol\"\n"
      "Delete Edition \"Delta\" \"Bob Smith\" 2019\n"
      "Delete Edition \"Delta\" \"Bob Smith\" 2020\n"
      "Return \"Delta\" \"Bob Smith\" 2020 \"Carol\"\n"
      "Delete Edition \"Delta\" \"Bob Smith\" 2020\n"
      "Find Book \"Delta\" \"Bob Smith\"\n"
      "Find Author \"Bob Smith\"\n"
      "Borrow \"Missing\" \"Nobody\" 1 \"Ghost\"\n"
      "Return \"Missing\" \"Nobody\" 1 \"Ghost\"\n"
      "Find Book \"Missing\" \"Nobody\"\n"
      "Find Author \"Nobody\"\n"
      "Delete Edition \"Missing\" \"Nobody\" 1\n"
      "Insert \"Zeta\" \"Ann\" -1 1\n"
      "Borrow \"Zeta\" \"Ann\" -1 \"Carol\"\n"
      "Find Book \"Zeta\" \"Ann\"\n"
      "Return \"Zeta\" \"Ann\" -1 \"Carol\"\n"
      "Delete Edition \"Zeta\" \"Ann\" -1\n"
      "Find Author \"Ann\"\n"
      "sort by Title\n"
      "Insert \"Lonely\"\n"
      "Borrow \"Alpha\" \"Zed\" 10\n"
      "Return \"Alpha\" \"Zed\" 10\n"
      "Find Book \"Alpha\"\n"
      "Find Author\n"
      "Delete Edition \"Alpha\" \"Zed\"\n"
      "Sort by Title extra\n"
      "Sort by Author\n");
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
      "Insert Zed's Alpha, Edition: 10.;"
      "Insert Ann's Alpha, Edition: 2.;"
      "Insert Ann's Alpha, Edition: 1.;"
      "Insert Ann's Beta, Edition: 3.;"
      "Insert Ann's Gamma Ray, Edition: 1.;"
      "Insert Bob Smith's Delta, Edition: 2020.;"
      "Insert Bob Smith's Delta, Edition: 2019.;"
      "Insert Zed's Omega, Edition: 0.;"
      "Datum already exist.;"
      "Title: Alpha	Author: Ann	Edition: <1(1), 2(2)>;"
      "Ann's Books: <Alpha, Beta, Gamma Ray>;"
      "Title: Alpha	Author: Ann	Edition: <1(1), 2(2)>;"
      "Title: Alpha	Author: Zed	Edition: <10(1)>;"
      "Title: Beta	Author: Ann	Edition: <3(1)>;"
      "Title: Delta	Author: Bob Smith	Edition: <2019(1), 2020(2)>;"
      "Title: Gamma Ray	Author: Ann	Edition: <1(3)>;"
      "Title: Omega	Author: Zed	Edition: <0(1)>;"
      "Title: Alpha	Author: Ann	Edition: <1(1), 2(2)>;"
      "Title: Beta	Author: Ann	Edition: <3(1)>;"
      "Title: Gamma Ray	Author: Ann	Edition: <1(3)>;"
      "Title: Delta	Author: Bob Smith	Edition: <2019(1), 2020(2)>;"
      "Title: Alpha	Author: Zed	Edition: <10(1)>;"
      "Title: Omega	Author: Zed	Edition: <0(1)>;"
      "Carol borrows Ann's Alpha, Edition: 1.;"
      "No available copies.;"
      "Already borrowed.;"
      "Outstanding borrows exist.;"
      "Carol borrows Ann's Alpha, Edition: 2.;"
      "Dave borrows Ann's Alpha, Edition: 2.;"
      "No available copies.;"
      "Title: Alpha	Author: Ann	Edition: <1(0), 2(0)>;"
      "Carol returns Ann's Alpha, Edition: 1.;"
      "Delete Ann's Alpha, Edition: 1.;"
      "Title: Alpha	Author: Ann	Edition: <2(0)>;"
      "Datum doesn't exist.;"
      "Borrow record not found.;"
      "Carol returns Ann's Alpha, Edition: 2.;"
      "Outstanding borrows exist.;"
      "Dave returns Ann's Alpha, Edition: 2.;"
      "Delete Ann's Alpha, Edition: 2.;"
      "Book doesn't exist.;"
      "Ann's Books: <Beta, Gamma Ray>;"
      "Carol borrows Bob Smith's Delta, Edition: 2019.;"
      "Carol borrows Bob Smith's Delta, Edition: 2020.;"
      "Title: Delta	Author: Bob Smith	Edition: <2019(0), 2020(1)>;"
      "Title: Alpha	Author: Zed	Edition: <10(1)>;"
      "Title: Beta	Author: Ann	Edition: <3(1)>;"
      "Title: Delta	Author: Bob Smith	Edition: <2019(0), 2020(1)>;"
      "Title: Gamma Ray	Author: Ann	Edition: <1(3)>;"
      "Title: Omega	Author: Zed	Edition: <0(1)>;"
      "Carol returns Bob Smith's Delta, Edition: 2019.;"
      "Delete Bob Smith's Delta, Edition: 2019.;"
      "Outstanding borrows exist.;"
      "Carol returns Bob Smith's Delta, Edition: 2020.;"
      "Delete Bob Smith's Delta, Edition: 2020.;"
      "Book doesn't exist.;"
      "No book found.;"
      "Datum doesn't exist.;"
      "Borrow record not found.;"
      "Book doesn't exist.;"
      "No book found.;"
      "Datum doesn't exist.;"
      "Insert Ann's Zeta, Edition: -1.;"
      "Carol borrows Ann's Zeta, Edition: -1.;"
      "Title: Zeta	Author: Ann	Edition: <-1(0)>;"
      "Carol returns Ann's Zeta, Edition: -1.;"
      "Delete Ann's Zeta, Edition: -1.;"
      "Ann's Books: <Beta, Gamma Ray>;"
      "Unknown Command.;"
      "Incomplete Command.;"
      "Incomplete Command.;"
      "Incomplete Command.;"
      "Incomplete Command.;"
      "Incomplete Command.;"
      "Incomplete Command.;"
      "Unknown Command.;"
      "Title: Beta	Author: Ann	Edition: <3(1)>;"
      "Title: Gamma Ray	Author: Ann	Edition: <1(3)>;"
      "Title: Alpha	Author: Zed	Edition: <10(1)>;"
      "Title: Omega	Author: Zed	Edition: <0(1)>;";

  std::cout.rdbuf(cout_buf);
  std::cin.rdbuf(cin_buf);
  CHECK(trimEnd(oss) == trimEnd(expected));
}

#endif