#pragma once
// Account.h
#ifndef ACCOUNT_H
#define ACCOUNT_H

#include <iostream>
using namespace std;

class NegativeDeposit {};
class InvalidAmount {};
class InsufficientFunds {};
class NegativeWithdrawal {};

class Account
{
private:
    double balance;

    bool isInteger(double amount)
    {
        return amount == (int)amount;
    }

public:
    Account()
    {
        balance = 0;
    }

    Account(double initialDeposit)
    {
        balance = initialDeposit;
    }

    double getBalance()
    {
        return balance;
    }

    double deposit(double amount)
    {
        if (!isInteger(amount))
            throw InvalidAmount();

        if (amount < 0)
            throw NegativeDeposit();

        balance += amount;
        return balance;
    }

    double withdraw(double amount)
    {
        if (!isInteger(amount))
            throw InvalidAmount();

        if (amount < 0)
            throw NegativeWithdrawal();

        if (amount > balance)
            throw InsufficientFunds();

        balance -= amount;
        return balance;
    }
};

#endif