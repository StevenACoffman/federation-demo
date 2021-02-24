## Apollo Federation Demo

This repository is a demo of using Apollo Federation to build a single schema on top of multiple services. The microservices are located under the [`./services`](./services/) folder and the gateway that composes the overall schema is in the [`gateway.js`](./gateway.js) file.

### Installation

To run this demo locally, pull down the repository then run the following commands:

```sh
npm install
```
This will install all of the dependencies for the gateway and each underlying service.

### Running Locally ( Gateway + Go)

There's a handy Go `main.go` script in the root that will run all the Go services and the Node.JS Apollo Gateway
```sh
go run main.go
```
| Port | Service |
|---|---|
| [4000](http://localhost:4000) | Gateway |
| [4001](http://localhost:4001)  | Accounts  |
| [4002](http://localhost:4002) | Reviews |
| [4003](http://localhost:4003)  | Products |
| [4004](http://localhost:4004) | Inventory |

### Running Locally (Node JS)

```sh
npm run start-services
```

This command will run all of the microservices at once. They can be found at http://localhost:4001, http://localhost:4002, http://localhost:4003, and http://localhost:4004.

In another terminal window, run the gateway by running this command:

```sh
npm run start-gateway
```

This will start up the gateway and serve it at http://localhost:4000

### What is this?

This demo showcases four partial schemas running as federated microservices. Each of these schemas can be accessed on their own and form a partial shape of an overall schema. The gateway fetches the service capabilities from the running services to create an overall composed schema which can be queried. 

To see the query plan when running queries against the gateway, click on the `Query Plan` tab in the bottom right hand corner of [GraphQL Playground](http://localhost:4000)

To learn more about Apollo Federation, check out the [docs](https://www.apollographql.com/docs/apollo-server/federation/introduction)

### Federated Combined Schema and Queries
If you open a browser to http://localhost:4000 you should be able to run two different top level queries.

Run any of these example queries as [http://localhost:4000/playground](http://localhost:4000/playground)

```
query MyReviews{
  me {
    username
    reviews {
      body
      product {
        name
        upc
      }
    }
  }
}
```

```
query TopProducts{
  topProducts{
    name
    reviews{
      author{
        name
      }
    }
  }
}
```
The combined federated schema is this:

```graphql
type Product {
upc: String!
name: String
price: Int
weight: Int
reviews: [Review]
inStock: Boolean
shippingEstimate: Int
}

type Query {
me: User
topProducts(first: Int = 5): [Product]
}

type Review {
id: ID!
body: String
author: User
product: Product
}

type User {
id: ID!
name: String
username: String
reviews: [Review]
}
```
### 

### DGraph
https://dgraph.io/docs/get-started/

Apollo Federation is now merged and available in the master branch of Dgraph. This is available via Docker Hub using the `dgraph/dgraph:master` Docker image. If you want to use a stable image tag (the master image always updates to the latest master), you can use `dgraph/dgraph:3642fed5`.

You can read more about how to use Apollo Federation in the docs currently in the description of PR #7275. This adds support for the @key, @extends, and @external directives.

This is slated to be released in the official Dgraph v21.03 version in March. Please do let us know if something doesn't work for you or if there's anything we can improve.

This PR extends support for the `Apollo Federation`.

## Support for Apollo federation directives
Our current implementation allows support for 3 directives, namely @key, @extends and @external.

### @key directive.
This directive is used on any type and it takes one field argument inside it which is called @key field. There are some limitations on how to use @key directives.

* User can define @key directive only once for a type, Support for multiple key types is not provided yet.
* Since the @key field act as a foreign key to resolve entities from the service where it is extended, the field provided as an argument inside @key directive should be of `ID` type or having `@id` directive on it. For example:-

```
type User @key(fields: "id") {
   id: ID!
  name: String
}
```

### @extends directive.
@extends directive is provided to give support for extended definitions. Suppose the above defined `User` type is defined in some service. Users can extend it to our GraphQL service by using this keyword.

```
type User @key(fields: "id") @extends{
   id: ID! @external
  products: [Product]
}
```

### @external directive.
@external directive means that the given field is not stored on this service. It is stored in some other service. This keyword can only be used on extended type definitions. Like it is used above on the `id` field.

## Generated Queries and mutations
In this section, we will mention what all the queries and mutations will be available to individual service and to the apollo `gateway`. We will take the given schema as our example:-

```
type Mission @key(fields: "id") {
    id: ID!
    crew: [Astronaut]
    designation: String!
    startDate: String
    endDate: String
}

type Astronaut @key(fields: "id") @extends {
    id: ID! @external
    missions: [Mission]
}
```

The queries and mutations which are exposed to the gateway are:-

```
type Query {
	getMission(id: ID!): Mission
	queryMission(filter: MissionFilter, order: MissionOrder, first: Int, offset: Int): [Mission]
	aggregateMission(filter: MissionFilter): MissionAggregateResult
}

type Mutation {
	addMission(input: [AddMissionInput!]!): AddMissionPayload
	updateMission(input: UpdateMissionInput!): UpdateMissionPayload
	deleteMission(filter: MissionFilter!): DeleteMissionPayload
	addAstronaut(input: [AddAstronautInput!]!): AddAstronautPayload
	updateAstronaut(input: UpdateAstronautInput!): UpdateAstronautPayload
	deleteAstronaut(filter: AstronautFilter!): DeleteAstronautPayload
}
```

The queries for `Astronaut` are not exposed to the gateway since it will be resolved through the `_entities` resolver. Although these queries will be available on the Dgraph GraphQL endpoint.

# Mutation for `extended` types
if we want to add an object of Astronaut type which is @Extended in this service.
The mutation `addAstronaut` takes `AddAstronautInput` which is generated as:-

```
input AddAstronautInput {
	id: ID!
	missions: [MissionRef]
}
```

Even though the `id` field is of `ID` type which should be ideally be generated internally by Dgraph, In this case, it should be provided as input. The reason for this is because of the unavailability of `federated mutations`, the user should provide the value of `id` same as the value present in the GraphQL service where the type `Astronaut` is defined.
For example, let's take that the type Astronaut is defined in some other service `AstronautService` as:-

```
type Astronaut @key(fields: "id") {
    id: ID! 
    name: String!
}
```

When adding an Object of `Astronaut` type, first it should be added into `AstronautService` and then the `addAstronaut` mutation should be called and value of `id` provided as an argument must be equal to the value in `AstronautService`.

# Gateway Supported Directives.
Due to the bug in the federation library (see [here](https://github.com/apollographql/federation/issues/346)), some directives are removed from the schema `SDL` which is returned to the gateway in response to the `_service` query.
Those directives are `@custom`, `@generate`, and `@auth`.
You can still use these directives in your GraphQL schema and they will work as desired but the gateway will unaware of this.

This change is

