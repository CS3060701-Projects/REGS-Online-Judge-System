#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "@SECRET@"

#include <sstream>
#include <string>

#include "Account.h"
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
  std::ostringstream oss;
  std::streambuf* cout_buf = std::cout.rdbuf();
  std::cout.rdbuf(oss.rdbuf());
  CHECK(true);

  Account a(100);
  try {
    cout << "Depositing 40" << ';';
    cout << "New balance: ";
    cout << a.deposit(40) << ';';
    cout << "Withdraw 15" << ';';
    cout << "New balance: ";
    cout << a.withdraw(15) << ';';
    cout << "Withdraw 350" << ';';
    cout << "New balance: ";
    cout << a.withdraw(350) << ';';
  } catch (InsufficientFunds) {
    cout << "Not enough money to withdraw that amount." << ';';
  } catch (NegativeDeposit) {
    cout << "You may only deposit a positive amount." << ';';
  }

  auto expected =
      "Depositing 40;"
      "New balance: 140;"
      "Withdraw 15;"
      "New balance: 125;"
      "Withdraw 350;"
      "New balance: Not enough money to withdraw that amount.;";

  std::cout.rdbuf(cout_buf);
  CHECK(trimEnd(oss) == trimEnd(expected));
}

#endif