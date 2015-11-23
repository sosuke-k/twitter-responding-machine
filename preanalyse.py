#! /usr/bin/env python
# -*- coding: utf-8 -*-

import os
import math
import nltk
import itertools
import matplotlib.pyplot as plt
from matplotlib.font_manager import FontProperties
from storm.locals import *
from mdcorpus.orm import *

# for Mac
font_path = '/Library/Fonts/Osaka.ttf'
font_prop = FontProperties(fname=font_path)

db_path = "dataset/corpus.db"

store = Store(create_database("sqlite:" + db_path))
lines = store.find(MovieLine)

print "Found %d MovieLine." % lines.count()
texts = [line.text for line in lines]

print "sentence tokenizing..."
sentences = itertools.chain(*[nltk.sent_tokenize(text.lower()) for text in texts])

print "word tokenizing..."
tokenized_sentences = [nltk.word_tokenize(sent) for sent in sentences]
print "Found %d sentences." % len(tokenized_sentences)

word_freq = nltk.FreqDist(itertools.chain(*tokenized_sentences))
print "Found %d unique words tokens." % len(word_freq.items())

items = sorted(word_freq.items(), key=lambda x: x[1], reverse=True)

print "Counting..."
dic = dict([(k, [x[0] for x in v])
            for (k, v) in itertools.groupby(items, key=lambda x: x[1])])
items = sorted(dic.items(), key=lambda x: x[0])

data = ""
x = []
y = []
for item in items:
    text = "出現回数%6d 回の単語%5d 個" % (item[0], len(item[1]))
    text = text + os.linesep
    sample = item[1][:10] if len(item[1]) >= 10 else item[1]
    print sample
    sample = [word.encode('utf-8') for word in sample]
    text = text + "  単語の例: " + ", ".join(sample)
    text = text + os.linesep
    data = data + text
    x.append(item[0])
    y.append(len(item[1]))

plt.bar(x, y)
plt.xscale('log')
plt.yscale('log')
plt.title(u"出現回数と単語数の関係", fontproperties=font_prop)
plt.xlabel(u"出現回数", fontproperties=font_prop)
plt.ylabel(u"単語数", fontproperties=font_prop)
plt.savefig("pre_analyse/tf_and_count.png")

with open('pre_analyse/tf_and_count.txt', 'w') as f:
    f.write(data)
