// @ts-ignore
export const apiLocation = import.meta.env.PROD ?
    window.location.origin :
    // `https://${window.location.hostname}:6325` :
    `https://${window.location.hostname}:6325`