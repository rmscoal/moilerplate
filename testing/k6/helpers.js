export function randomCharacters(length) {
  const charset = "abcdefghijklmnopqrstuvwxyz";
  let res = "";
  while (length--) res += charset[(Math.random() * charset.length) | 0];
  return res;
}

/**
 * Generate random string of digits
 * @param {number} length
 * @returns {string}
 */
export function randomNumbers(length) {
  const charset = "123456789";
  let res = "";
  while (length--) res += charset[(Math.random() * charset.length) | 0];
  return res;
}

export function generateNewUser() {
  return {
    name: randomCharacters(10),
    email: `${randomCharacters(10)}@gmail.com`,
    username: randomCharacters(20),
    password: "verystrongpassword",
    phoneNumber: `${randomNumbers(16)}`,
  }
}

