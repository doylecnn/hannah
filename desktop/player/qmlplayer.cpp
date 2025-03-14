#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QScreen>
#include <QTimer>

#include "qmlplayer.h"
#include "FlacPic.h"
#include "ID3v2Pic.h"
#include "bassplayer.h"
#include "configurationwindow.h"
#include "lrcbar.h"
#include "lyrics.h"
#include "osd.h"
#include "playlistmanagewindow.h"

QmlPlayer::QmlPlayer(QObject *parent)
    : QObject(parent),
      m_timer(new QTimer()),
      m_lrcTimer(new QTimer()),
      m_lyrics(new Lyrics()),
      m_osd(new OSD()),
      m_lb(new LrcBar(m_lyrics, gBassPlayer))
{
    gBassPlayer->devInit();
    connect(m_timer, &QTimer::timeout, this, &QmlPlayer::onUpdateTime);
    connect(m_lrcTimer, &QTimer::timeout, this, &QmlPlayer::onUpdateLrc);

    m_timer->start(27);
    m_lrcTimer->start(70);

#if defined(Q_OS_WIN) && QT_VERSION < QT_VERSION_CHECK(6, 0, 0)
    taskbarButton   = new QWinTaskbarButton(this);
    taskbarProgress = taskbarButton->progress();
    taskbarProgress->setRange(0, 1000);

    thumbnailToolBar   = new QWinThumbnailToolBar(this);
    playToolButton     = new QWinThumbnailToolButton(thumbnailToolBar);
    stopToolButton     = new QWinThumbnailToolButton(thumbnailToolBar);
    backwardToolButton = new QWinThumbnailToolButton(thumbnailToolBar);
    forwardToolButton  = new QWinThumbnailToolButton(thumbnailToolBar);
    playToolButton->setToolTip(tr("Player"));
    playToolButton->setIcon(QIcon(":/rc/images/player/Play.png"));
    stopToolButton->setToolTip(tr("Stop"));
    stopToolButton->setIcon(QIcon(":/rc/images/player/Stop.png"));
    backwardToolButton->setToolTip(tr("Previous"));
    backwardToolButton->setIcon(QIcon(":/rc/images/player/Pre.png"));
    forwardToolButton->setToolTip(tr("Next"));
    forwardToolButton->setIcon(QIcon(":/rc/images/player/Next.png"));
    thumbnailToolBar->addButton(playToolButton);
    thumbnailToolBar->addButton(stopToolButton);
    thumbnailToolBar->addButton(backwardToolButton);
    thumbnailToolBar->addButton(forwardToolButton);
    connect(playToolButton, &QWinThumbnailToolButton::clicked, this, &QmlPlayer::onPlay);
    connect(stopToolButton, &QWinThumbnailToolButton::clicked, this, &QmlPlayer::onPlayStop);
    connect(backwardToolButton, &QWinThumbnailToolButton::clicked, this, &QmlPlayer::onPlayPrevious);
    connect(forwardToolButton, &QWinThumbnailToolButton::clicked, this, &QmlPlayer::onPlayNext);
#endif
}

QmlPlayer::~QmlPlayer()
{
    delete m_lb;
    delete m_osd;
    delete m_lyrics;
    delete m_lrcTimer;
    delete m_timer;
}

void QmlPlayer::setTaskbarButtonWindow()
{
#if defined(Q_OS_WIN) && QT_VERSION < QT_VERSION_CHECK(6, 0, 0)
    auto *window = qobject_cast<QWindow *>(gQmlApplicationEngine->rootObjects().at(0));
    taskbarButton->setWindow(window);
    thumbnailToolBar->setWindow(window);
#endif
}

void QmlPlayer::onQuit()
{
    QCoreApplication::quit();
}

void QmlPlayer::onShowPlaylists()
{
    Q_ASSERT(playlistManageWindow);
    if (playlistManageWindow->isHidden())
        playlistManageWindow->showNormal();
    playlistManageWindow->raise();
    playlistManageWindow->activateWindow();
}

void QmlPlayer::onSettings()
{
    Q_ASSERT(configurationWindow);
    configurationWindow->onShowConfiguration();
}

void QmlPlayer::onFilter() {}

void QmlPlayer::onMessage() {}

void QmlPlayer::onMusic() {}

void QmlPlayer::onCloud() {}

void QmlPlayer::onBluetooth() {}

void QmlPlayer::onCart() {}

void QmlPlayer::presetEQChanged(int index)
{
    QVector<QVector<int>> presets = {{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
                                     {3, 1, 0, -2, -4, -4, -2, 0, 1, 2},
                                     {-2, 0, 2, 4, -2, -2, 0, 0, 4, 4},
                                     {-6, 1, 4, -2, -2, -4, 0, 0, 6, 6},
                                     {0, 8, 8, 4, 0, 0, 0, 0, 2, 2},
                                     {-6, 0, 0, 0, 0, 0, 4, 0, 4, 0},
                                     {-2, 3, 4, 1, -2, -2, 0, 0, 4, 4},
                                     {-2, 0, 0, 2, 2, 0, 0, 0, 4, 4},
                                     {0, 0, 0, 4, 4, 4, 0, 2, 3, 4},
                                     {-2, 0, 2, 1, 0, 0, 0, 0, -2, -4},
                                     {-4, 0, 2, 1, 0, 0, 0, 0, -4, -6},
                                     {0, 0, 0, 4, 5, 3, 6, 3, 0, 0},
                                     {-4, 0, 2, 0, 0, 0, 0, 0, -4, -6}};
    if (index >= 0 && index < presets.length())
    {
        auto &preset = presets[index];
        setEq0(preset[0]);
        setEq1(preset[1]);
        setEq2(preset[2]);
        setEq3(preset[3]);
        setEq4(preset[4]);
        setEq5(preset[5]);
        setEq6(preset[6]);
        setEq7(preset[7]);
        setEq8(preset[8]);
        setEq9(preset[9]);
        emit eq0Changed();
        emit eq1Changed();
        emit eq2Changed();
        emit eq3Changed();
        emit eq4Changed();
        emit eq5Changed();
        emit eq6Changed();
        emit eq7Changed();
        emit eq8Changed();
        emit eq9Changed();
    }
}

void QmlPlayer::onOpenPreset() {}

void QmlPlayer::onSavePreset() {}

void QmlPlayer::onFavorite() {}

void QmlPlayer::onStop() {}

void QmlPlayer::onPrevious() {}

void QmlPlayer::onPause() {}

void QmlPlayer::onNext() {}

void QmlPlayer::onRepeat() {}

void QmlPlayer::onShuffle() {}

void QmlPlayer::onSwitchFiles() {}

void QmlPlayer::onSwitchPlaylists() {}

void QmlPlayer::onSwitchFavourites() {}

void QmlPlayer::onOpenFile() {}

int QmlPlayer::getPrimaryScreenWidth()
{
    auto *screen = QGuiApplication::primaryScreen();
    Q_ASSERT(screen);
    auto size = screen->availableSize();
    return size.width();
}

int QmlPlayer::getPrimaryScreenHeight()
{
    auto *screen = QGuiApplication::primaryScreen();
    Q_ASSERT(screen);
    auto size = screen->availableSize();
    return size.height();
}

void QmlPlayer::showNormal()
{
    emit showPlayer();
}

void QmlPlayer::loadAudio(const QString &uri) {}

void QmlPlayer::addToListAndPlay(const QList<QUrl> &uris) {}

void QmlPlayer::addToListAndPlay(const QStringList &uris) {}

void QmlPlayer::addToListAndPlay(const QString &uri) {}

qreal QmlPlayer::getEq0() const
{
    return m_eq0;
}

qreal QmlPlayer::getEq1() const
{
    return m_eq1;
}

qreal QmlPlayer::getEq2() const
{
    return m_eq2;
}

qreal QmlPlayer::getEq3() const
{
    return m_eq3;
}

qreal QmlPlayer::getEq4() const
{
    return m_eq4;
}

qreal QmlPlayer::getEq5() const
{
    return m_eq5;
}

qreal QmlPlayer::getEq6() const
{
    return m_eq6;
}

qreal QmlPlayer::getEq7() const
{
    return m_eq7;
}

qreal QmlPlayer::getEq8() const
{
    return m_eq8;
}

qreal QmlPlayer::getEq9() const
{
    return m_eq9;
}

qreal QmlPlayer::getVolumn() const
{
    return m_volumn;
}

qreal QmlPlayer::getProgress() const
{
    return m_progress;
}

const QString &QmlPlayer::getCoverUrl() const
{
    return m_coverUrl;
}

const QString &QmlPlayer::getSongName() const
{
    return m_songName;
}

void QmlPlayer::setEq0(qreal value)
{
    if (m_eq0 == value)
        return;
    m_eq0 = value;
    Q_ASSERT(gBassPlayer);
    gBassPlayer->setEQ(0, (int)m_eq0);
}

void QmlPlayer::setEq1(qreal value)
{
    if (m_eq1 == value)
        return;
    m_eq1 = value;
    Q_ASSERT(gBassPlayer);
    gBassPlayer->setEQ(1, (int)m_eq1);
}

void QmlPlayer::setEq2(qreal value)
{
    if (m_eq2 == value)
        return;
    m_eq2 = value;
    Q_ASSERT(gBassPlayer);
    gBassPlayer->setEQ(2, (int)m_eq2);
}

void QmlPlayer::setEq3(qreal value)
{
    if (m_eq3 == value)
        return;
    m_eq3 = value;
    Q_ASSERT(gBassPlayer);
    gBassPlayer->setEQ(3, (int)m_eq3);
}

void QmlPlayer::setEq4(qreal value)
{
    if (m_eq4 == value)
        return;
    m_eq4 = value;
    Q_ASSERT(gBassPlayer);
    gBassPlayer->setEQ(4, (int)m_eq4);
}

void QmlPlayer::setEq5(qreal value)
{
    if (m_eq5 == value)
        return;
    m_eq5 = value;
    Q_ASSERT(gBassPlayer);
    gBassPlayer->setEQ(5, (int)m_eq6);
}

void QmlPlayer::setEq6(qreal value)
{
    if (m_eq6 == value)
        return;
    m_eq6 = value;
    Q_ASSERT(gBassPlayer);
    gBassPlayer->setEQ(6, (int)m_eq6);
}

void QmlPlayer::setEq7(qreal value)
{
    if (m_eq7 == value)
        return;
    m_eq7 = value;
    Q_ASSERT(gBassPlayer);
    gBassPlayer->setEQ(7, (int)m_eq7);
}

void QmlPlayer::setEq8(qreal value)
{
    if (m_eq8 == value)
        return;
    m_eq8 = value;
    Q_ASSERT(gBassPlayer);
    gBassPlayer->setEQ(8, (int)m_eq8);
}

void QmlPlayer::setEq9(qreal value)
{
    if (m_eq9 == value)
        return;
    m_eq9 = value;
    Q_ASSERT(gBassPlayer);
    gBassPlayer->setEQ(9, (int)m_eq9);
}

void QmlPlayer::setVolumn(qreal value)
{
    if (m_volumn == value)
        return;
    m_volumn = value;
    Q_ASSERT(gBassPlayer);
    gBassPlayer->setVol((int)m_volumn);
}

void QmlPlayer::setProgress(qreal progress)
{
    if (m_progress == progress)
        return;
    m_progress = progress;
}

void QmlPlayer::setCoverUrl(const QString &u)
{
    if (m_coverUrl == u)
        return;
    m_coverUrl = u;
}

void QmlPlayer::setSongName(const QString &n)
{
    if (m_songName == n)
        return;
    m_songName = n;
}

void QmlPlayer::onPlay() {}

void QmlPlayer::onPlayStop() {}

void QmlPlayer::onPlayPrevious() {}

void QmlPlayer::onPlayNext() {}

void QmlPlayer::onUpdateTime() {}

void QmlPlayer::onUpdateLrc() {}
