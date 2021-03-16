//
// C++ Implementation: parsercsv
//
// Description: module to parse Comma-Separated Values (CSV) formatted playlists (rfc4180)
//
//
// Author: Ingo Kossyk <kossyki@cs.tu-berlin.de>, (C) 2004
// Author: Tobias Rafreider trafreider@mixxx.org, (C) 2011
// Author: Daniel Schürmann daschuer@gmx.de, (C) 2011
//
// Copyright: See COPYING file that comes with this distribution
//
//

#include <QDir>
#include <QMessageBox>
#include <QTextStream>
#include <QtDebug>

#include "parsercsv.h"
#include "moc_parsercsv.cpp"

QList<QString> ParserCsv::parse(const QString &sFilename)
{
    QFile   file(sFilename);
    QString basepath = sFilename.section('/', 0, -2);

    clearLocations();
    // qDebug() << "ParserCsv: Starting to parse.";
    if (file.open(QIODevice::ReadOnly) && !isBinary(sFilename))
    {
        QByteArray ba = file.readAll();

        QList<QList<QString>> tokens = tokenize(ba, ',');

        // detect Location column
        int loc_coll = 0x7fffffff;
        if (tokens.size())
        {
            for (int i = 0; i < tokens[0].size(); ++i)
            {
                if (tokens[0][i] == tr("Location"))
                {
                    loc_coll = i;
                    break;
                }
            }
            for (int i = 1; i < tokens.size(); ++i)
            {
                if (loc_coll < tokens[i].size())
                {
                    // Todo: check if path is relative
                    QFileInfo fi = tokens[i][loc_coll];
                    if (fi.isRelative())
                    {
                        // add base path
                        qDebug() << "is relative" << basepath << fi.filePath();
                        fi.setFile(basepath, fi.filePath());
                    }
                    m_sLocations.append(fi.filePath());
                }
            }
        }

        file.close();

        if (m_sLocations.count() != 0)
        {
            return m_sLocations;
        }
        else
        {
            return QList<QString>(); // NULL pointer returned when no locations were found
        }
    }

    file.close();
    return QList<QString>(); // if we get here something went wrong
}

// Code was posted at http://www.qtcentre.org/threads/35511-Parsing-CSV-data
// by "adzajac" and adapted to use QT Classes
QList<QList<QString>> ParserCsv::tokenize(const QByteArray &str, char delimiter)
{
    QList<QList<QString>> tokens;

    unsigned int row    = 0;
    bool         quotes = false;
    QByteArray   field  = "";

    tokens.append(QList<QString>());

    for (int pos = 0; pos < str.length(); ++pos)
    {
        char c = str[pos];
        if (!quotes && c == '"')
        {
            quotes = true;
        }
        else if (quotes && c == '"')
        {
            if (pos + 1 < str.length() && str[pos + 1] == '"')
            {
                field.append(c);
                pos++;
            }
            else
            {
                quotes = false;
            }
        }
        else if (!quotes && c == delimiter)
        {
            if (isUtf8(field.constData()))
            {
                tokens[row].append(QString::fromUtf8(field));
            }
            else
            {
                tokens[row].append(QString::fromLatin1(field));
            }
            field.clear();
        }
        else if (!quotes && (c == '\r' || c == '\n'))
        {
            if (isUtf8(field.constData()))
            {
                tokens[row].append(QString::fromUtf8(field));
            }
            else
            {
                tokens[row].append(QString::fromLatin1(field));
            }
            field.clear();
            tokens.append(QList<QString>());
            row++;
        }
        else
        {
            field.push_back(c);
        }
    }
    return tokens;
}
