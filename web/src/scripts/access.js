export function accessible(need, have) {
  switch (need) {
    case "administrator":
      switch (have) {
        case "administrator":
          return true;
        default:
          return false;
      }
    case "operator":
      switch (have) {
        case "administrator":
          return true;
        case "operator":
          return true;
        default:
          return false;
      }
    default:
      return true;
  }
}
