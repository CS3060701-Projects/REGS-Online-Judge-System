#ifndef ACCOUNT_H
#define ACCOUNT_H

class Account
{
private:
	double balance;
public:
	Account() { balance = 0; }
	Account(double initialDeposit) { balance = initialDeposit; }
	double getBalance() { return balance; }


	// returns new balance or -1 if error
	double deposit(double amount)
	{
		if (amount > 0) balance += amount;
		else return -1;
		return balance;
	}
	// return new balance or -1 if invalid amount
	double withdraw(double amount)
	{
		if ((amount > balance) || (amount < 0)) return -1;
		else balance -= amount;
		return balance;
	}
};

#endif