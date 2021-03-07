import {useState} from "react";
import useFetch from "react-fetch-hook";

export function useServer(setData, setList, list) {
  const [pkgUrl, setPkgUrl] = useState('')
  const { isLoading, data, error } = useFetch(pkgUrl, { depends: [pkgUrl], mode: "no-cors"});
  console.log('DATA:', data)
  console.log('isLoading', isLoading)
  if (!error) {
    console.log('no error')
  }
  return setPkgUrl
}
