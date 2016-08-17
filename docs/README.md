# Passions Hacked - Search Engine

We have available a database of:

| City  | Passion | Endorsement Count |
| ------------- | ------------- | ------------- | 
| Amsterdam  | Museums | 64380 |
| Amsterdam  | City Walks | 41419 |
| Amsterdam  | Sightseeing | 37198 |
| Amsterdam  | Architecture | 33933 |
| Amsterdam  | Walking | 32707 |
| Amsterdam  | Culture | 32543 |
| Amsterdam  | Nightlife | 31500 |
| Amsterdam  | Shopping | 25756 |
| Amsterdam  | Cycling | 24234 |
| Lisboa  | Food | 33060 |
| Lisboa  | City Walks | 32724 |
| Lisboa  | Old Town | 31719 |
| Lisboa  | Sightseeing | 30543 |
| Lisboa  | Friendly People | 26728 |
| Lisboa  | Architecture | 26207 |
| Lisboa  | Culture | 26070 |
| Lisboa  | History | 26009 |
| Lisboa  | Monuments | 25904 |
| Paris  | Museums | 84569 |
| Paris  | Sightseeing | 63457 |
| Paris  | Culture | 58852 |
| Paris  | Monuments | 55475 |
| Paris  | Architecture | 54666 |
| Paris  | Shopping | 47174 |
| Paris  | History | 47054 |
| Paris  | Food | 45947 |
| Paris  | Art | 42927 |
| ... | ... | ... |

And given a list of passions we want to find out destinations around the
world that matches those passions.

### How would I implement a search engine to give the best recommendations to our users

#### First try (the dumb one)

Since we don't have much data to begin with, my first attempt would use
a simple strategy. A high level description of the approach would be:

* have the data indexed by passion
* for each city keep a sum of all endorsements it received and the 
maximum endorsements a passion on that city has
* for each passion on the passions list received retrieve a list of 
cities that are endorsed for that passion
* based on the lists retrieved compute a score for each city
* sort cities by score

##### Score Calculation

The initial score for each city is composed from two factors:

* `passion_relevance`: sum of endorsements counts normalized according to
all cities endorsed for some of the users' passions
 - for each search passion `p` we have a list `Lp` of cities that were
endorsed for that passion, with that we can find out `max_endorsement_p`
which is the maximum endorsement any city has for passion `p`. Then, if 
`endorsement_count_p` is the amount of endorsements that a city has for 
passion `p`, the `passion_relevance` of a city is the sum, from all query
passions `p` it has associated, of `endorsement_count_p` / `max_endorsement_p`.
 - `passion_relevance` is then a number between `0` and `number of user
passions` and the bigger it is the more relevant the city is for the 
user passions. It is a crucial info since cities with higher values are
the ones most endorsed for most of users' passions.

* `city_passions_relevance`: sum of endorsements counts normalized
according to all passions a city is endorsed for
 - again, `endorsement_count_p` is the amount of endorsements that a 
city has for passion `p`. If `max_city_endorsement` is the maximum 
endorsements a passion on that city has, the `city_passions_relevance` 
of a city is the sum, for all user passions `p` it has associated, of
`endorsement_count_p` / `max_city_endorsement`. 
 - `city_passions_relevance` is also a number between `0` and `number of
user passions` and the bigger it is the more relevant the user passions
are for the city. It's an important info to use since it tells us if
users' passions are what the city is most known for.

Given these, the score of a city is simply the sum of `passion_relevance`
+ 0.6 * `city_passions_relevance`.
We multiply `city_passions_relevance` by 0.6 since we believe 
`passion_relevance` is more decisive.


#### Next Steps (making our search smarter)

---

First step to improve our search engine is change the way passion 
matching works, for now only exact matches are taken into account, we 
should start using some analysis on the query tokens ASAP to display
correct results despite typos.

---

- With the data we have currently there's no much room for improving our
search results, we could go some steps further and use more advanced
techniques but the precision improvement probably wouldn't be that
expressive. So the first steps I would take are:
  1. create a validation dataset!
  2. look for all available data related to this task I could find to
create more factors the score should account for.

To make sure we are really increasing the accuracy of our search engine
we **must** have a validation dataset. We must manually define which
cities we would like to have associated with a query or come up with
ways of algorithmically doing this. With such dataset in hands we can
now be sure that further attempts on improving the search engine are
actually improving it.

Also, counts of endorsements a city has for a passion are barely enough
for finding places where the user will fully experience its passions.
To have results that will blow our users mind we need more info about
the user, and users with similar passions. With data from places the
users traveled and whether they liked it or not we might start finding
correlations and figuring out places where people with similar passions
went for enjoying those passion.

---

- Now that we have more data and ways to validate our improvements its
time to start using advanced techniques and machine learning in our aid.
  1. [Tie-Yan Liu. **Learning to Rank for Information Retrieval**](http://dl.acm.org/citation.cfm?id=1618304)
looks like a good read to start with.

Once we have more factors to take into account when calculating the score
of each city we might start using machine learning techniques for figuring
out the weights each factor should have. Also, find correlations between
destinations and passions, and use these as new factors on score calculation. 

### Questions remaining

#### Popularity bias?

What's popular grows in popularity by default! Despite being familiar
with the concept I haven't studied any techniques for overcoming it yet.
I believe that the `city_passions_relevance` factor used on the score
calculation helps mitigate its effects but is not enough. I believe it's
strongly related to finding hidden gems, since popular destinations are
probably popular on all passions it is endorsed for, but with the info
we currently have I'm not confident to increase it's weight. It's safer
to keep popular places growing in popularity than shows bad results to
users. Once we have more data available a natural step would be to
normalize factors according to popularity but probably the machine
learning techniques we would be using at this point would already take
care of this for us.


#### Ignoring jokes?

If users start endorsing destinations as jokes it would be a hard
problem for us to figure out without more data. I believe that having
more datasets we could find these jokes more easily, and allowing users
to report these kind of stuff would also help. Anyway, I like to believe
most of our users would be traveling and living their passions instead
of doing these kind of jokes, and that our ranking algorithms are good
enough that this jokes wouldn't have a big impact on search quality,
so it's one of the last things I would care about.

#### What else do we need to think about?

- metrics, metrics and more metrics.. click rate, conversion rate and
everything else we  can think about. we need ways of validating our
search results.
- creating a system for human evaluation of search results would
certainly be a good idea


# That' all folks!
