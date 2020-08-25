export default {
  methods: {
    graphQLErrorCode(error, code) {
      if (!error.graphQLErrors) return false;
      for (const ge of error.graphQLErrors) {
        if (ge.extensions?.code == code) return true;
      }
      return false;
    },
    graphQLErrors(error) {
      // eslint-disable-next-line no-console
      if (!error.graphQLErrors) console.log("error", error);
      return error.graphQLErrors
        .reduce((ax, err) => {
          ax.push(err.message);
          return ax;
        }, [])
        .join(", ");
    },
  },
};
