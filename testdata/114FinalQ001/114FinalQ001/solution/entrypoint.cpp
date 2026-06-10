#define CATCH_CONFIG_RUNNER
#define CATCH_CONFIG_DELIM "\t"
#define ONLINE_JUDGE
#include "case.h"

int main(int argc, char* argv[]) {
  Catch::Session session;

  auto& config = session.configData();
  config.showSuccessfulTests = false;
  config.reporterName = "compact";
  config.abortAfter = -1;
  config.useColour = Catch::UseColour::No;

  auto failedCount = session.run(argc, argv);

  return failedCount;
}