const roles = {
  administrator: 2,
  operator: 1,
  developer: 0,
};

export function accessible(need, have) {
  return roles[have] >= roles[need];
}
