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

TEST_CASE("case2") {
  std::ostringstream oss;
  std::streambuf* cout_buf = std::cout.rdbuf();
  std::cout.rdbuf(oss.rdbuf());

  Account a(350);

  try {
    cout << "Depositing -50" << ';';
    cout << "New balance: ";
    cout << a.deposit(-50) << ';';
  } catch (NegativeDeposit) {
    cout << "You may only deposit a positive amount." << ';';
  }

  try {
    cout << "Depositing 30.5" << ';';
    cout << "New balance: ";
    cout << a.deposit(30.5) << ';';
  } catch (InvalidAmount) {
    cout << "Amount must be a non-negative integer." << ';';
  }
  CHECK(true);

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
  try {
    cout << "Withdraw -50" << ';';
    cout << "New balance: ";
    cout << a.withdraw(-50) << ';';
  } catch (NegativeWithdrawal) {
    cout << "Withdrawal amount cannot be negative." << ';';
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
      "New balance: Withdrawal amount cannot be negative.;";
  std::cout.rdbuf(cout_buf);
  CHECK(trimEnd(oss) == trimEnd(expected));
}

#endif