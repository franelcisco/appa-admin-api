package shopify

// getCustomerPartnerID is the GraphQL query to get a customer's metafield
const getCustomerPartnerID = `
query CustomerPartnerID($id: ID!, $key: String!, $ns: String!) {
  customer(id: $id) {
    id
	displayName
    metafield(namespace: $ns, key: $key) {
      id
      namespace
      key
      value
    }
  }
}`

// markOrderAsPaid is the GraphQL mutation to mark an order as paid
const markOrderAsPaid = `
mutation orderMarkAsPaid($id: ID!) {
  orderMarkAsPaid(input: { id: $id })
  {
    order {
      id
    }
    userErrors {
      field
      message
    }
  }
}
`
