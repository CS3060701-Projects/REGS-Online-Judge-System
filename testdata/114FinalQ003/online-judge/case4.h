#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "1f4ad5e0f89587eb282ac0f00dd68dd32b26cdaab264496652b1ece1805d7009"

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

TEST_CASE("case4") {
  std::ostringstream oss;
  std::streambuf* cout_buf = std::cout.rdbuf();
  std::cout.rdbuf(oss.rdbuf());
  CHECK(true);

  Account a(350);

  try {
    cout << "Depositing -50" << ';';
    cout << "New balance: ";
    cout << a.deposit(-50) << ';';
  } catch (NegativeDeposit) {
    cout << "You may only deposit a positive amount." << ';';
  }
  CHECK(true);
  CHECK(true);

  try {
    cout << "Depositing 30.5" << ';';
    cout << "New balance: ";
    cout << a.deposit(30.5) << ';';
  } catch (InvalidAmount) {
    cout << "Amount must be a non-negative integer." << ';';
  }

  try {
    cout << "Withdraw 250" << ';';
    cout << "New balance: ";
    cout << a.withdraw(250) << ';';
    cout << "Withdraw 70" << ';';
    cout << "New balance: ";
    cout << a.withdraw(70) << ';';
    cout << "Withdraw 50" << ';';
    cout << "New balance: ";
    cout << a.withdraw(50) << ';';
  } catch (InsufficientFunds) {
    cout << "Not enough money to withdraw that amount." << ';';
  }
  CHECK(true);
  CHECK(true);
  CHECK(true);

  try {
    cout << "Withdraw -50" << ';';
    cout << "New balance: ";
    cout << a.withdraw(-50) << ';';
  } catch (NegativeWithdrawal) {
    cout << "Withdrawal amount cannot be negative." << ';';
  }

  try {
    cout << "Withdraw 0.5" << ';';
    cout << "New balance: ";
    cout << a.withdraw(0.5) << ';';
  } catch (InvalidAmount) {
    cout << "Amount must be a non-negative integer." << ';';
  }

  auto expected =
      "Depositing -50;"
      "New balance: You may only deposit a positive amount.;"
      "Depositing 30.5;"
      "New balance: Amount must be a non-negative integer.;"
      "Withdraw 250;"
      "New balance: 100;"
      "Withdraw 70;"
      "New balance: 30;"
      "Withdraw 50;"
      "New balance: Not enough money to withdraw that amount.;"
      "Withdraw -50;"
      "New balance: Withdrawal amount cannot be negative.;"
      "Withdraw 0.5;"
      "New balance: Amount must be a non-negative integer.;";
  std::cout.rdbuf(cout_buf);
  CHECK(trimEnd(oss) == trimEnd(expected));
}

#endif
