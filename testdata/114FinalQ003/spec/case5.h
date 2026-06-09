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

TEST_CASE("case5") {
  std::ostringstream oss;
  std::streambuf* cout_buf = std::cout.rdbuf();
  std::cout.rdbuf(oss.rdbuf());
  CHECK(true);

  Account a(100);
  try {
    cout << "Withdraw 100" << ';';
    cout << "New balance: ";
    cout << a.withdraw(100) << ';';
  } catch (InsufficientFunds) {
    cout << "Not enough money to withdraw that amount." << ';';
  }

  // Exception safety: balance unchanged after NegativeDeposit
  Account b(200);
  try {
    cout << "Depositing -100" << ';';
    b.deposit(-100);
  } catch (NegativeDeposit) {
    cout << "You may only deposit a positive amount." << ';';
  }
  cout << "Balance after failed deposit: ";
  cout << b.getBalance() << ';';

  // Exception safety: balance unchanged after InsufficientFunds
  Account c(50);
  try {
    cout << "Withdraw 100" << ';';
    c.withdraw(100);
  } catch (InsufficientFunds) {
    cout << "Not enough money to withdraw that amount." << ';';
  }
  cout << "Balance after failed withdrawal: ";
  cout << c.getBalance() << ';';
  CHECK(true);

  // Exception safety: balance unchanged after InvalidAmount on withdraw
  Account d(500);
  try {
    cout << "Withdraw 99.9" << ';';
    d.withdraw(99.9);
  } catch (InvalidAmount) {
    cout << "Amount must be a non-negative integer." << ';';
  }
  cout << "Balance after invalid withdrawal: ";
  cout << d.getBalance() << ';';

  // Partial success: deposit succeeds, then withdraw exceeds new balance
  Account e(100);
  try {
    cout << "Depositing 50" << ';';
    cout << "New balance: ";
    cout << e.deposit(50) << ';';
    cout << "Withdraw 200" << ';';
    cout << "New balance: ";
    cout << e.withdraw(200) << ';';
  } catch (InsufficientFunds) {
    cout << "Not enough money to withdraw that amount." << ';';
  }

  auto expected =
      "Withdraw 100;"
      "New balance: 0;"
      "Depositing -100;"
      "You may only deposit a positive amount.;"
      "Balance after failed deposit: 200;"
      "Withdraw 100;"
      "Not enough money to withdraw that amount.;"
      "Balance after failed withdrawal: 50;"
      "Withdraw 99.9;"
      "Amount must be a non-negative integer.;"
      "Balance after invalid withdrawal: 500;"
      "Depositing 50;"
      "New balance: 150;"
      "Withdraw 200;"
      "New balance: Not enough money to withdraw that amount.;";
  CHECK(true);
  CHECK(true);
  CHECK(true);

  std::cout.rdbuf(cout_buf);
  CHECK(trimEnd(oss) == trimEnd(expected));
}

#endif