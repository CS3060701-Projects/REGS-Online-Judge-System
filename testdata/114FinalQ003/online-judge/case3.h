#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "77361bc1a7f896ba8d5810b6a1b0231088191bb4547464b278e0616b32be9147"

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

TEST_CASE("case3") {
  std::ostringstream oss;
  std::streambuf* cout_buf = std::cout.rdbuf();
  std::cout.rdbuf(oss.rdbuf());
  CHECK(true);

  Account a(200);

  try {
    cout << "Depositing 75" << ';';
    cout << "New balance: ";
    cout << a.deposit(75) << ';';
    cout << "Withdraw 55" << ';';
    cout << "New balance: ";
    cout << a.withdraw(55) << ';';
    cout << "Depositing -25" << ';';
    cout << "New balance: ";
    cout << a.deposit(-25) << ';';
  } catch (NegativeDeposit) {
    cout << "You may only deposit a positive amount." << ';';
  }

  try {
    cout << "Withdraw 250" << ';';
    cout << "New balance: ";
    cout << a.withdraw(250) << ';';
  } catch (InsufficientFunds) {
    cout << "Not enough money to withdraw that amount." << ';';
  }
  CHECK(true);

  try {
    cout << "Withdraw -50" << ';';
    cout << "New balance: ";
    cout << a.withdraw(-50) << ';';
  } catch (NegativeWithdrawal) {
    cout << "Withdrawal amount cannot be negative." << ';';
  }
  CHECK(true);
  CHECK(true);

  auto expected =
      "Depositing 75;"
      "New balance: 275;"
      "Withdraw 55;"
      "New balance: 220;"
      "Depositing -25;"
      "New balance: You may only deposit a positive amount.;"
      "Withdraw 250;"
      "New balance: Not enough money to withdraw that amount.;"
      "Withdraw -50;"
      "New balance: Withdrawal amount cannot be negative.;";
  std::cout.rdbuf(cout_buf);
  CHECK(trimEnd(oss) == trimEnd(expected));
}

#endif
