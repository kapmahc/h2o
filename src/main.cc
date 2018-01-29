#include <Wt/WApplication.h>
#include <Wt/WBreak.h>
#include <Wt/WContainerWidget.h>
#include <Wt/WLineEdit.h>
#include <Wt/WPushButton.h>
#include <Wt/WText.h>

class HelloApplication : public Wt::WApplication {
public:
  HelloApplication(const Wt::WEnvironment &env);

private:
  Wt::WLineEdit *nameEdit_;
  Wt::WText *greeting_;

  void greet();
};

HelloApplication::HelloApplication(const Wt::WEnvironment &env)
    : Wt::WApplication(env) {
  setTitle("Hello world");

  root()->addWidget(std::make_unique<Wt::WText>("Your name, please? "));
  nameEdit_ = root()->addWidget(std::make_unique<Wt::WLineEdit>());
  Wt::WPushButton *button =
      root()->addWidget(std::make_unique<Wt::WPushButton>("Greet me."));
  root()->addWidget(std::make_unique<Wt::WBreak>());
  greeting_ = root()->addWidget(std::make_unique<Wt::WText>());
  button->clicked().connect(this, [this]() {
    greeting_->setText("Hello there, " + nameEdit_->text());
  });
}

int main(int argc, char **argv) {
  return Wt::WRun(argc, argv, [](const Wt::WEnvironment &env) {
    return std::make_unique<HelloApplication>(env);
  });
}
