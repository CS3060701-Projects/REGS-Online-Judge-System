#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "1a23b2755b59d90b1e88a1d60c1f2ef7e3b72a8c367805ab083c2db686aa3be4"

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

TEST_CASE("case3") {
  std::istringstream inputCommand(
      "Insert \"Zen\" \"Alan Watts\" 1 2\n"
      "Insert \"Zen\" \"Alan Watts\" 2 1\n"
      "Insert \"Zen\" \"Shunryu Suzuki\" 1 3\n"
      "Insert \"Art of War\" \"Sun Tzu\" 1 1\n"
      "Insert \"Art of War\" \"Sun Tzu\" 2 2\n"
      "Insert \"Meditations\" \"Marcus Aurelius\" 1 4\n"
      "Insert \"Meditations\" \"Marcus Aurelius\" 2 2\n"
      "Insert \"Meditations\" \"Marcus Aurelius\" 3 1\n"
      "Find Author \"Alan Watts\"\n"
      "Find Author \"Sun Tzu\"\n"
      "Sort by Title\n"
      "Borrow \"Zen\" \"Alan Watts\" 1 \"Mike\"\n"
      "Borrow \"Zen\" \"Alan Watts\" 2 \"Mike\"\n"
      "Borrow \"Zen\" \"Shunryu Suzuki\" 1 \"Mike\"\n"
      "Borrow \"Art of War\" \"Sun Tzu\" 1 \"Mike\"\n"
      "Borrow \"Art of War\" \"Sun Tzu\" 1 \"Nancy\"\n"
      "Borrow \"Zen\" \"Alan Watts\" 1 \"Mike\"\n"
      "Borrow \"Art of War\" \"Sun Tzu\" 99 \"Mike\"\n"
      "Borrow \"Meditations\" \"Marcus Aurelius\" 1 \"Alice\"\n"
      "Borrow \"Meditations\" \"Marcus Aurelius\" 2 \"Alice\"\n"
      "Delete Edition \"Meditations\" \"Marcus Aurelius\" 1\n"
      "Delete Edition \"Meditations\" \"Marcus Aurelius\" 2\n"
      "Delete Edition \"Meditations\" \"Marcus Aurelius\" 3\n"
      "Find Book \"Meditations\" \"Marcus Aurelius\"\n"
      "Delete Edition \"Art of War\" \"Sun Tzu\" 1\n"
      "Return \"Art of War\" \"Sun Tzu\" 1 \"Mike\"\n"
      "Delete Edition \"Art of War\" \"Sun Tzu\" 1\n"
      "Delete Edition \"Art of War\" \"Sun Tzu\" 5\n"
      "Find Book \"Art of War\" \"Sun Tzu\"\n"
      "Delete Edition \"Zen\" \"Alan Watts\" 1\n"
      "Delete Edition \"Zen\" \"Alan Watts\" 2\n"
      "Find Book \"Zen\" \"Alan Watts\"\n"
      "Borrow \"Zen\" \"Alan Watts\" 1 \"Nancy\"\n"
      "Return \"Zen\" \"Alan Watts\" 1 \"Mike\"\n"
      "Delete Edition \"Zen\" \"Alan Watts\" 1\n"
      "Delete Edition \"Zen\" \"Alan Watts\" 2\n"
      "Find Author \"Alan Watts\"\n"
      "Return \"Zen\" \"Alan Watts\" 2 \"Mike\"\n"
      "Return \"Zen\" \"Alan Watts\" 2 \"Nobody\"\n"
      "Borrow \"Meditations\" \"Marcus Aurelius\" 1 \"Bob\"\n"
      "Borrow \"Meditations\" \"Marcus Aurelius\" 2 \"Bob\"\n"
      "Delete Edition \"Meditations\" \"Marcus Aurelius\" 1\n"
      "Delete Edition \"Meditations\" \"Marcus Aurelius\" 2\n"
      "Sort by Author\n"
      "Return \"Meditations\" \"Marcus Aurelius\" 1 \"Alice\"\n"
      "Return \"Meditations\" \"Marcus Aurelius\" 2 \"Alice\"\n"
      "Sort by Title\n");
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
      "Insert Alan Watts's Zen, Edition: 1.;"
      "Insert Alan Watts's Zen, Edition: 2.;"
      "Insert Shunryu Suzuki's Zen, Edition: 1.;"
      "Insert Sun Tzu's Art of War, Edition: 1.;"
      "Insert Sun Tzu's Art of War, Edition: 2.;"
      "Insert Marcus Aurelius's Meditations, Edition: 1.;"
      "Insert Marcus Aurelius's Meditations, Edition: 2.;"
      "Insert Marcus Aurelius's Meditations, Edition: 3.;"
      "Alan Watts's Books: <Zen>;"
      "Sun Tzu's Books: <Art of War>;"
      "Title: Art of War	Author: Sun Tzu	Edition: <1(1), 2(2)>;"
      "Title: Meditations	Author: Marcus Aurelius	Edition: <1(4), 2(2), "
      "3(1)>;"
      "Title: Zen	Author: Alan Watts	Edition: <1(2), 2(1)>;"
      "Title: Zen	Author: Shunryu Suzuki	Edition: <1(3)>;"
      "Mike borrows Alan Watts's Zen, Edition: 1.;"
      "Mike borrows Alan Watts's Zen, Edition: 2.;"
      "Mike borrows Shunryu Suzuki's Zen, Edition: 1.;"
      "Mike borrows Sun Tzu's Art of War, Edition: 1.;"
      "No available copies.;"
      "Already borrowed.;"
      "Datum doesn't exist.;"
      "Alice borrows Marcus Aurelius's Meditations, Edition: 1.;"
      "Alice borrows Marcus Aurelius's Meditations, Edition: 2.;"
      "Outstanding borrows exist.;"
      "Outstanding borrows exist.;"
      "Delete Marcus Aurelius's Meditations, Edition: 3.;"
      "Title: Meditations	Author: Marcus Aurelius	Edition: <1(3), 2(1)>;"
      "Outstanding borrows exist.;"
      "Mike returns Sun Tzu's Art of War, Edition: 1.;"
      "Delete Sun Tzu's Art of War, Edition: 1.;"
      "Datum doesn't exist.;"
      "Title: Art of War	Author: Sun Tzu	Edition: <2(2)>;"
      "Outstanding borrows exist.;"
      "Outstanding borrows exist.;"
      "Title: Zen	Author: Alan Watts	Edition: <1(1), 2(0)>;"
      "Nancy borrows Alan Watts's Zen, Edition: 1.;"
      "Mike returns Alan Watts's Zen, Edition: 1.;"
      "Outstanding borrows exist.;"
      "Outstanding borrows exist.;"
      "Alan Watts's Books: <Zen>;"
      "Mike returns Alan Watts's Zen, Edition: 2.;"
      "Borrow record not found.;"
      "Bob borrows Marcus Aurelius's Meditations, Edition: 1.;"
      "Bob borrows Marcus Aurelius's Meditations, Edition: 2.;"
      "Outstanding borrows exist.;"
      "Outstanding borrows exist.;"
      "Title: Zen	Author: Alan Watts	Edition: <1(1), 2(1)>;"
      "Title: Meditations	Author: Marcus Aurelius	Edition: <1(2), 2(0)>;"
      "Title: Zen	Author: Shunryu Suzuki	Edition: <1(2)>;"
      "Title: Art of War	Author: Sun Tzu	Edition: <2(2)>;"
      "Alice returns Marcus Aurelius's Meditations, Edition: 1.;"
      "Alice returns Marcus Aurelius's Meditations, Edition: 2.;"
      "Title: Art of War	Author: Sun Tzu	Edition: <2(2)>;"
      "Title: Meditations	Author: Marcus Aurelius	Edition: <1(3), 2(1)>;"
      "Title: Zen	Author: Alan Watts	Edition: <1(1), 2(1)>;"
      "Title: Zen	Author: Shunryu Suzuki	Edition: <1(2)>;";

  std::cout.rdbuf(cout_buf);
  std::cin.rdbuf(cin_buf);
  CHECK(trimEnd(oss) == trimEnd(expected));
}

#endif
