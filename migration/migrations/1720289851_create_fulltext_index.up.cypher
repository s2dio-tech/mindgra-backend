CREATE FULLTEXT INDEX contentAndDescriptions FOR (n:Word) ON EACH [n.content, n.description]
OPTIONS {
  indexConfig: {
    `fulltext.analyzer`: 'english',
    `fulltext.eventually_consistent`: true
  }
}