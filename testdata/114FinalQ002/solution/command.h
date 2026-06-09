#ifndef _COMMAND_H_
#define _COMMAND_H_

#include <algorithm>
#include <iomanip>
#include <iostream>
#include <map>
#include <set>
#include <sstream>
#include <string>
#include <vector>

using namespace std;

struct EditionInfo {
  int quantity;
  set<string> borrowers;
};

struct Book {
  string title;
  string author;
  map<int, EditionInfo> editions;
};

map<pair<string, string>, Book> libraryDB;

inline bool hasExtraToken(stringstream& ss) {
  string extra;
  if (ss >> extra) {
    return true;
  }
  return false;
}

inline void printBook(const Book& b) {
  cout << "Title: " << b.title << "\tAuthor: " << b.author << "\tEdition: <";

  bool first = true;

  for (map<int, EditionInfo>::const_iterator it = b.editions.begin();
       it != b.editions.end(); ++it) {
    if (!first) {
      cout << ", ";
    }

    cout << it->first << "(" << it->second.quantity << ")";
    first = false;
  }

  cout << ">" << ';';
}

inline void execute(std::string command) {
  {
    if (command.empty() || command.find_first_not_of(" \r\n\t") == string::npos) {
      return;
    }

    stringstream ss(command);
    string cmd;
    ss >> cmd;

    // Insert "title" "author" edition quantity
    if (cmd == "Insert") {
      string title, author;
      int edition, quantity;

      if (!(ss >> quoted(title) >> quoted(author) >> edition >> quantity)) {
        cout << "Incomplete Command." << ';';
        return;
      }

      if (hasExtraToken(ss)) {
        cout << "Unknown Command." << ';';
        return;
      }

      pair<string, string> key = make_pair(title, author);

      if (libraryDB.count(key) && libraryDB[key].editions.count(edition)) {
        cout << "Datum already exist." << ';';
        return;
      }

      libraryDB[key].title = title;
      libraryDB[key].author = author;
      libraryDB[key].editions[edition].quantity = quantity;

      cout << "Insert " << author << "'s " << title << ", Edition: " << edition
           << "." << ';';
    }

    // Delete Edition "title" "author" edition
    else if (cmd == "Delete") {
      string type;

      if (!(ss >> type)) {
        cout << "Incomplete Command." << ';';
        return;
      }

      if (type != "Edition") {
        cout << "Unknown Command." << ';';
        return;
      }

      string title, author;
      int edition;

      if (!(ss >> quoted(title) >> quoted(author) >> edition)) {
        cout << "Incomplete Command." << ';';
        return;
      }

      if (hasExtraToken(ss)) {
        cout << "Unknown Command." << ';';
        return;
      }

      pair<string, string> key = make_pair(title, author);

      if (!libraryDB.count(key) || !libraryDB[key].editions.count(edition)) {
        cout << "Datum doesn't exist." << ';';
        return;
      }

      if (!libraryDB[key].editions[edition].borrowers.empty()) {
        cout << "Outstanding borrows exist." << ';';
        return;
      }

      libraryDB[key].editions.erase(edition);

      cout << "Delete " << author << "'s " << title << ", Edition: " << edition
           << "." << ';';

      if (libraryDB[key].editions.empty()) {
        libraryDB.erase(key);
      }
    }

    // Find Book "title" "author"
    // Find Author "author"
    else if (cmd == "Find") {
      string type;

      if (!(ss >> type)) {
        cout << "Incomplete Command." << ';';
        return;
      }

      if (type == "Book") {
        string title, author;

        if (!(ss >> quoted(title) >> quoted(author))) {
          cout << "Incomplete Command." << ';';
          return;
        }

        if (hasExtraToken(ss)) {
          cout << "Unknown Command." << ';';
          return;
        }

        pair<string, string> key = make_pair(title, author);

        if (!libraryDB.count(key)) {
          cout << "Book doesn't exist." << ';';
          return;
        }

        printBook(libraryDB[key]);
      } else if (type == "Author") {
        string author;

        if (!(ss >> quoted(author))) {
          cout << "Incomplete Command." << ';';
          return;
        }

        if (hasExtraToken(ss)) {
          cout << "Unknown Command." << ';';
          return;
        }

        vector<string> titles;

        for (map<pair<string, string>, Book>::const_iterator it =
                 libraryDB.begin();
             it != libraryDB.end(); ++it) {
          if (it->second.author == author) {
            titles.push_back(it->second.title);
          }
        }

        sort(titles.begin(), titles.end());
        titles.erase(unique(titles.begin(), titles.end()), titles.end());

        if (titles.empty()) {
          cout << "No book found." << ';';
          return;
        }

        cout << author << "'s Books: <";

        for (int i = 0; i < (int)titles.size(); ++i) {
          if (i) {
            cout << ", ";
          }

          cout << titles[i];
        }

        cout << ">" << ';';
      } else {
        cout << "Unknown Command." << ';';
      }
    }

    // Sort by Title
    // Sort by Author
    else if (cmd == "Sort") {
      string by, type;

      if (!(ss >> by >> type)) {
        cout << "Incomplete Command." << ';';
        return;
      }

      if (by != "by") {
        cout << "Unknown Command." << ';';
        return;
      }

      if (hasExtraToken(ss)) {
        cout << "Unknown Command." << ';';
        return;
      }

      vector<Book> books;

      for (map<pair<string, string>, Book>::const_iterator it =
               libraryDB.begin();
           it != libraryDB.end(); ++it) {
        books.push_back(it->second);
      }

      if (type == "Title") {
        sort(books.begin(), books.end(), [](const Book& a, const Book& b) {
          if (a.title != b.title) {
            return a.title < b.title;
          }

          return a.author < b.author;
        });
      } else if (type == "Author") {
        sort(books.begin(), books.end(), [](const Book& a, const Book& b) {
          if (a.author != b.author) {
            return a.author < b.author;
          }

          return a.title < b.title;
        });
      } else {
        cout << "Unknown Command." << ';';
        return;
      }

      for (int i = 0; i < (int)books.size(); ++i) {
        printBook(books[i]);
      }
    }

    // Borrow "title" "author" edition "borrower"
    else if (cmd == "Borrow") {
      string title, author, borrower;
      int edition;

      if (!(ss >> quoted(title) >> quoted(author) >> edition >>
            quoted(borrower))) {
        cout << "Incomplete Command." << ';';
        return;
      }

      if (hasExtraToken(ss)) {
        cout << "Unknown Command." << ';';
        return;
      }

      pair<string, string> key = make_pair(title, author);

      if (!libraryDB.count(key) || !libraryDB[key].editions.count(edition)) {
        cout << "Datum doesn't exist." << ';';
        return;
      }

      EditionInfo& info = libraryDB[key].editions[edition];

      if (info.borrowers.count(borrower)) {
        cout << "Already borrowed." << ';';
        return;
      }

      if (info.quantity <= 0) {
        cout << "No available copies." << ';';
        return;
      }

      --info.quantity;
      info.borrowers.insert(borrower);

      cout << borrower << " borrows " << author << "'s " << title
           << ", Edition: " << edition << "." << ';';
    }

    // Return "title" "author" edition "borrower"
    else if (cmd == "Return") {
      string title, author, borrower;
      int edition;

      if (!(ss >> quoted(title) >> quoted(author) >> edition >>
            quoted(borrower))) {
        cout << "Incomplete Command." << ';';
        return;
      }

      if (hasExtraToken(ss)) {
        cout << "Unknown Command." << ';';
        return;
      }

      pair<string, string> key = make_pair(title, author);

      if (!libraryDB.count(key) || !libraryDB[key].editions.count(edition) ||
          !libraryDB[key].editions[edition].borrowers.count(borrower)) {
        cout << "Borrow record not found." << ';';
        return;
      }

      EditionInfo& info = libraryDB[key].editions[edition];

      ++info.quantity;
      info.borrowers.erase(borrower);

      cout << borrower << " returns " << author << "'s " << title
           << ", Edition: " << edition << "." << ';';
    }

    else {
      cout << "Unknown Command." << ';';
    }
  }
}

#endif