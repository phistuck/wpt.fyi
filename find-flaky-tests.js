const URL = require('url').URL;
const fetch = require('node-fetch');
const cTable = require('console.table');

const numberOfRuns = 10;
const browsers = ['chrome', 'firefox'];

async function findFlakes(browser) {
  let url = new URL('/api/runs', 'https://staging.wpt.fyi');
  url.searchParams.set('labels', 'master,experimental');
  url.searchParams.set('product', browser);
  url.searchParams.set('max-count', 10);
  const ids = await fetch(url)
    .then(r => r.json())
    .then(runs => runs.map(r => r.id));
  const query = {
    "run_ids": ids,
    "query":{
      "and": [
        {
          "or": [
            {"browser_name": browser, "status":"PASS"},
            {"browser_name": browser, "status":"OK"}
          ]
        },
        {
          "or": [
            {"browser_name": browser, "status":"TIMEOUT"},
            {"browser_name": browser, "status":"ERROR"},
            {"browser_name": browser, "status":"FAIL"}
          ]
        }
      ]
    }
  };
  return fetch('https://staging.wpt.fyi/api/search', {
    method: 'POST',
    body: JSON.stringify(query),
  }).then(
    async r => {
      if (!r.ok) {
        console.log(await r.text());
        throw 'Failed to fetch';
      }
      return r.json();
    }
  );
}

function getAllFlakes() {
  const flakes = new Map();
  Promise.all(
    browsers.map(b => findFlakes(b))
  ).then(searches => {
    const browserRates = searches.map(s => {
      return s.results.map(r => {
        const passes = r.legacy_status.reduce((sum, next) => sum + next.passes / next.total, 0);
        const rate = Math.round(passes * 100) / 100 / numberOfRuns;
        return {
          'test': r.test,
          'passRate': rate,
        };
      });
    });
    for (const i in browserRates) {
      const rates = browserRates[i];
      for (const rate of rates) {
        if (!flakes.has(rate.test)) {
          flakes.set(rate.test, browsers.map(b => 0));
        }
        flakes.get(rate.test)[i] = rate.passRate;
      }
    }
    const table = [['Test', ...browsers]];
    for (const entry of flakes.entries()) {
      table.push([entry[0], ...entry[1]]);
    }
    console.table(table);
  });
}

getAllFlakes();