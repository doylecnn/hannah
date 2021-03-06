#include <QApplication>
#include <QCommandLineParser>
#include <QDir>
#include <QFileOpenEvent>
#include <QMessageBox>
#include <QSettings>
#include <QtCore>

#include "configurationwindow.h"

#if defined(Q_OS_MACOS)
#    include "application.h"
#    include "serviceslots.h"
#else
#    if defined(Q_OS_WIN)
#        include <Windows.h>
#        include <shellapi.h>
#        include <tchar.h>
#    endif
#    include "qtsingleapplication.h"
#endif

#if defined(Q_OS_MACOS)

static ConfigurationWindow *cw = nullptr;

void serviceSearch(const QString &s)
{
    if (cw)
    {
        cw->onMacServiceSearch(s);
    }
}
void serviceOpenUrl(const QString &s)
{
    if (cw)
    {
        cw->onMacServiceSearch(s);
    }
}

void serviceOpenLink(const QString &s)
{
    if (cw)
    {
        cw->onMacServiceOpenLink(s);
    }
}

void serviceAppendToPlaylist(const QStringList &s)
{
    if (cw)
    {
        cw->onMacServiceAppendToPlaylist(s);
    }
}

void serviceClearAndAddToPlaylist(const QStringList &s)
{
    if (cw)
    {
        cw->onMacServiceClearAndAddToPlaylist(s);
    }
}

void serviceAppendToPlaylistFile(const QStringList &s)
{
    if (cw)
    {
        cw->onMacServiceAppendToPlaylistFile(s);
    }
}

void serviceClearAndAddToPlaylistFile(const QStringList &s)
{
    if (cw)
    {
        cw->onMacServiceClearAndAddToPlaylistFile(s);
    }
}

#endif

void i18n(QTranslator &translator, QTranslator &qtTranslator)
{
    QString locale = "zh_CN";

    // main application and dynamic linked library locale
    QString localeDirPath = QCoreApplication::applicationDirPath() +
#if defined(Q_OS_MACOS)
                            "/../Resources/translations";
#else
                            "/translations";
#endif

    if (translator.load("Hannah_" + locale, localeDirPath))
    {
        qDebug() << "loading Hannah" << locale << " from " << localeDirPath << " success";
        if (QCoreApplication::installTranslator(&translator))
        {
            qDebug() << "installing translator success ";
        }
    }

    if (qtTranslator.load("qt_" + locale, localeDirPath))
    {
        qDebug() << "loading qt" << locale << " from " << localeDirPath << " success";
        if (QCoreApplication::installTranslator(&qtTranslator))
        {
            qDebug() << "installing qt translator success ";
        }
    }
}

int main(int argc, char *argv[])
{
    QTranslator translator;
    QTranslator qtTranslator;
#if defined(Q_OS_MACOS)
    Application a(argc, argv);
    i18n(translator, qtTranslator);
    ConfigurationWindow  w;
    w.connect(&a, &Application::openUrl, &w, &ConfigurationWindow::onOpenUrl);

    cw = &w;

    void registerHannahService();
    registerHannahService();
#else
    QtSingleApplication a(argc, argv);
    QCoreApplication::setApplicationName("Hannah");
    QCoreApplication::setApplicationVersion("1.0");

    QCommandLineParser parser;
    parser.setApplicationDescription("Hannah");
    parser.addHelpOption();
    parser.addVersionOption();

    parser.process(a);

    const QStringList args = parser.positionalArguments();

    if (a.isRunning())
    {
        if (args.length() > 0)
        {
            a.sendMessage(args.join("~"));
        }
        return 0;
    }

    i18n(translator, qtTranslator);
    ConfigurationWindow w;
    w.connect(&a, &QtSingleApplication::messageReceived, &w, &ConfigurationWindow::onApplicationMessageReceived);
    if (args.length() > 0)
    {
        w.onApplicationMessageReceived(args.join("~"));
    }
#    if defined(Q_OS_WIN)
    else
    {
        QSettings mxKey("HKEY_CLASSES_ROOT\\hannah", QSettings::NativeFormat);
        QString   v1 = mxKey.value(".").toString();
        QSettings mxOpenKey("HKEY_CLASSES_ROOT\\hannah\\shell\\open\\command", QSettings::NativeFormat);
        QString   v2 = mxOpenKey.value(".").toString();

        if (v1 != "URL:hannah Protocol" ||
            v2 != QChar('"') + QDir::toNativeSeparators(QCoreApplication::applicationFilePath()) + QString("\" \"%1\""))
        {
            QString cmd = QDir::toNativeSeparators(QCoreApplication::applicationDirPath() + "/registerProtocolHandler.exe");
            ::ShellExecuteW(nullptr,
                            L"open",
                            cmd.toStdWString().c_str(),
                            nullptr,
                            QDir::toNativeSeparators(QCoreApplication::applicationDirPath()).toStdWString().c_str(),
                            SW_SHOWNORMAL);
        }
    }
#    endif
#endif

    a.setQuitOnLastWindowClosed(false);
    return a.exec();
}
