package graphiql

import (
	"html/template"
)

type dataTmpl struct {
	Endpoint          string
	Es6PromiseVersion string
	FetchVersion      string
	ReactVersion      string
	GraphiQLVersion   string
}

const (
	templateName      = "graphiql"
	es6PromiseVersion = "4.2.6"
	fetchVersion      = "2.0.4"
	reactVersion      = "16.8.4"
	GqlVersion        = "0.13.0" // version of graphiql js library
)

func data(endpoint string) *dataTmpl {
	return &dataTmpl{
		Endpoint:          endpoint,
		Es6PromiseVersion: es6PromiseVersion,
		FetchVersion:      fetchVersion,
		ReactVersion:      reactVersion,
		GraphiQLVersion:   GqlVersion,
	}
}

func preparingTemplate() (*template.Template, error) {
	t := template.New(templateName)
	t, err := t.Parse(rawTemplate)
	if err != nil {
		return nil, err
	}

	return t, nil
}

var rawTemplate = `
<!--
 *  Copyright (c) Facebook, Inc.
 *  All rights reserved.
 *
 *  This source code is licensed under the license found in the
 *  LICENSE file in the root directory of this source tree.
-->
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            height: 100%;
            margin: 0;
            width: 100%;
            overflow: hidden;
        }
        #graphiql {
            height: 100vh;
        }
    </style>
    <!--
      This GraphiQL example depends on Promise and fetch, which are available in
      modern browsers, but can be "polyfilled" for older browsers.
      GraphiQL itself depends on React DOM.
      If you do not want to rely on a CDN, you can host these files locally or
      include them directly in your favored resource bunder.
    -->
	<script src="//cdn.jsdelivr.net/npm/es6-promise@{{ .Es6PromiseVersion }}/dist/es6-promise.auto.min.js"></script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/fetch/{{ .FetchVersion }}/fetch.min.js"></script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/react/{{ .ReactVersion }}/umd/react.production.min.js"></script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/react-dom/{{ .ReactVersion }}/umd/react-dom.production.min.js"></script>
    <!--
      These two files can be found in the npm module, however you may wish to
      copy them directly into your environment, or perhaps include them in your
      favored resource bundler.
     -->
    <link rel="stylesheet" href="//cdn.jsdelivr.net/npm/graphiql@{{ .GraphiQLVersion }}/graphiql.css" />
    <script src="//cdn.jsdelivr.net/npm/graphiql@{{ .GraphiQLVersion }}/graphiql.js"></script>
</head>
<body>
<div id="graphiql">Loading...</div>
<script>
    /**
     * This GraphiQL example illustrates how to use some of GraphiQL's props
     * in order to enable reading and updating the URL parameters, making
     * link sharing of queries a little bit easier.
     *
     * This is only one example of this kind of feature, GraphiQL exposes
     * various React params to enable interesting integrations.
     */
            // Parse the search string to get url parameters.
    var search = window.location.search;
    var parameters = {};
    search.substr(1).split('&').forEach(function (entry) {
        var eq = entry.indexOf('=');
        if (eq >= 0) {
            parameters[decodeURIComponent(entry.slice(0, eq))] =
                    decodeURIComponent(entry.slice(eq + 1));
        }
    });
    // if variables was provided, try to format it.
    if (parameters.variables) {
        try {
            parameters.variables =
                    JSON.stringify(JSON.parse(parameters.variables), null, 2);
        } catch (e) {
            // Do nothing, we want to display the invalid JSON as a string, rather
            // than present an error.
        }
    }
    // When the query and variables string is edited, update the URL bar so
    // that it can be easily shared
    function onEditQuery(newQuery) {
        parameters.query = newQuery;
        updateURL();
    }
    function onEditVariables(newVariables) {
        parameters.variables = newVariables;
        updateURL();
    }
    function onEditOperationName(newOperationName) {
        parameters.operationName = newOperationName;
        updateURL();
    }
    function updateURL() {
        var newSearch = '?' + Object.keys(parameters).filter(function (key) {
            return Boolean(parameters[key]);
        }).map(function (key) {
            return encodeURIComponent(key) + '=' +
                    encodeURIComponent(parameters[key]);
        }).join('&');
        history.replaceState(null, null, newSearch);
    }
    // Defines a GraphQL fetcher using the fetch API. You're not required to
    // use fetch, and could instead implement graphQLFetcher however you like,
    // as long as it returns a Promise or Observable.
    function graphQLFetcher(graphQLParams) {
        // This example expects a GraphQL server at the path /graphql.
        // Change this to point wherever you host your GraphQL server.
        return fetch('{{ .Endpoint }}', {
            method: 'post',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(graphQLParams),
        }).then(function (response) {
            return response.text();
        }).then(function (responseBody) {
            try {
                return JSON.parse(responseBody);
            } catch (error) {
                return responseBody;
            }
        });
    }
    // Render <GraphiQL /> into the body.
    // See the README in the top level of this module to learn more about
    // how you can customize GraphiQL by providing different values or
    // additional child elements.
    ReactDOM.render(
            React.createElement(GraphiQL, {
                fetcher: graphQLFetcher,
                query: parameters.query,
                variables: parameters.variables,
                operationName: parameters.operationName,
                onEditQuery: onEditQuery,
                onEditVariables: onEditVariables,
                onEditOperationName: onEditOperationName
            }),
            document.getElementById('graphiql')
    );
</script>
</body>
</html>
`
