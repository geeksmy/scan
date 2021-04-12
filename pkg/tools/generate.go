package tools

import (
	"fmt"
	"strings"
	"sync"
)

func Generate(beforeStr, rearStr, specialStr []string, l int, strs chan string, genWG *sync.WaitGroup) {
	defer genWG.Done()
	var wg sync.WaitGroup

	// 无特殊字符
	wg.Add(1)
	go func() {
		for i := 0; i < len(beforeStr); i++ {
			for j := 0; j < len(rearStr); j++ {
				strs <- fmt.Sprintf("%s%s", beforeStr[i], rearStr[j])
				strs <- fmt.Sprintf("%s%s", strings.Title(beforeStr[i]), strings.Title(rearStr[j]))
			}
		}
		wg.Done()
	}()

	// 有特殊字符 前
	wg.Add(1)
	go func() {
		for i := 0; i < len(beforeStr); i++ {
			for j := 0; j < len(rearStr); j++ {
				switch l {
				case 2:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							strs <- fmt.Sprintf("%s%s%s%s", specialStr[z], specialStr[x], beforeStr[i], rearStr[j])
							strs <- fmt.Sprintf("%s%s%s%s", specialStr[z], specialStr[x], strings.Title(beforeStr[i]), strings.Title(rearStr[j]))
						}
					}
				case 3:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							for c := 0; c < len(specialStr); c++ {
								strs <- fmt.Sprintf("%s%s%s%s%s", specialStr[z], specialStr[x], specialStr[c], beforeStr[i], rearStr[j])
								strs <- fmt.Sprintf("%s%s%s%s%s", specialStr[z], specialStr[x], specialStr[c], strings.Title(beforeStr[i]), strings.Title(rearStr[j]))
							}
						}
					}
				default:
					for z := 0; z < len(specialStr); z++ {
						strs <- fmt.Sprintf("%s%s%s", specialStr[z], beforeStr[i], rearStr[j])
						strs <- fmt.Sprintf("%s%s%s", specialStr[z], strings.Title(beforeStr[i]), strings.Title(rearStr[j]))
					}
				}
			}
		}
		wg.Done()
	}()

	// 有特殊字符 中
	wg.Add(1)
	go func() {
		for i := 0; i < len(beforeStr); i++ {
			for j := 0; j < len(rearStr); j++ {
				switch l {
				case 2:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							strs <- fmt.Sprintf("%s%s%s%s", beforeStr[i], specialStr[z], specialStr[x], rearStr[j])
							strs <- fmt.Sprintf("%s%s%s%s", strings.Title(beforeStr[i]), specialStr[z], specialStr[x], strings.Title(rearStr[j]))
						}
					}
				case 3:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							for c := 0; c < len(specialStr); c++ {
								strs <- fmt.Sprintf("%s%s%s%s%s", beforeStr[i], specialStr[z], specialStr[x], specialStr[c], rearStr[j])
								strs <- fmt.Sprintf("%s%s%s%s%s", strings.Title(beforeStr[i]), specialStr[z], specialStr[x], specialStr[c], strings.Title(rearStr[j]))
							}
						}
					}
				default:
					for z := 0; z < len(specialStr); z++ {
						strs <- fmt.Sprintf("%s%s%s", beforeStr[i], specialStr[z], rearStr[j])
						strs <- fmt.Sprintf("%s%s%s", strings.Title(beforeStr[i]), specialStr[z], strings.Title(rearStr[j]))
					}
				}
			}
		}
		wg.Done()
	}()

	// 有特殊字符 后
	wg.Add(1)
	go func() {
		for i := 0; i < len(beforeStr); i++ {
			for j := 0; j < len(rearStr); j++ {
				switch l {
				case 2:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							strs <- fmt.Sprintf("%s%s%s%s", beforeStr[i], rearStr[j], specialStr[z], specialStr[x])
							strs <- fmt.Sprintf("%s%s%s%s", strings.Title(beforeStr[i]), strings.Title(rearStr[j]), specialStr[z], specialStr[x])
						}
					}
				case 3:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							for c := 0; c < len(specialStr); c++ {
								strs <- fmt.Sprintf("%s%s%s%s%s", beforeStr[i], rearStr[j], specialStr[z], specialStr[x], specialStr[c])
								strs <- fmt.Sprintf("%s%s%s%s%s", strings.Title(beforeStr[i]), strings.Title(rearStr[j]), specialStr[z], specialStr[x], specialStr[c])
							}
						}
					}
				default:
					for z := 0; z < len(specialStr); z++ {
						strs <- fmt.Sprintf("%s%s%s", beforeStr[i], rearStr[j], specialStr[z])
						strs <- fmt.Sprintf("%s%s%s", strings.Title(beforeStr[i]), strings.Title(rearStr[j]), specialStr[z])
					}
				}
			}
		}
		wg.Done()
	}()

	// 有特殊字符 前中
	wg.Add(1)
	go func() {
		for i := 0; i < len(beforeStr); i++ {
			for j := 0; j < len(rearStr); j++ {
				switch l {
				case 2:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							for v := 0; v < len(specialStr); v++ {
								for b := 0; b < len(specialStr); b++ {
									strs <- fmt.Sprintf("%s%s%s%s%s%s", specialStr[z], specialStr[x], beforeStr[i],
										specialStr[v], specialStr[b], rearStr[j])
									strs <- fmt.Sprintf("%s%s%s%s%s%s", specialStr[z], specialStr[x], strings.Title(beforeStr[i]),
										specialStr[v], specialStr[b], strings.Title(rearStr[j]))
								}
							}
						}
					}
				case 3:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							for c := 0; c < len(specialStr); c++ {
								for v := 0; v < len(specialStr); v++ {
									for b := 0; b < len(specialStr); b++ {
										for n := 0; n < len(specialStr); n++ {
											strs <- fmt.Sprintf("%s%s%s%s%s%s%s%s", specialStr[z], specialStr[x], specialStr[c],
												beforeStr[i], specialStr[v], specialStr[b], specialStr[n], rearStr[j])
											strs <- fmt.Sprintf("%s%s%s%s%s%s%s%s", specialStr[z], specialStr[x], specialStr[c],
												strings.Title(beforeStr[i]), specialStr[v], specialStr[b], specialStr[n], strings.Title(rearStr[j]))
										}
									}
								}
							}
						}
					}
				default:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							strs <- fmt.Sprintf("%s%s%s%s", specialStr[z], beforeStr[i], specialStr[x], rearStr[j])
							strs <- fmt.Sprintf("%s%s%s%s", specialStr[z], strings.Title(beforeStr[i]), specialStr[x], strings.Title(rearStr[j]))
						}
					}
				}
			}
		}
		wg.Done()
	}()

	// 有特殊字符 中后
	wg.Add(1)
	go func() {
		for i := 0; i < len(beforeStr); i++ {
			for j := 0; j < len(rearStr); j++ {
				switch l {
				case 2:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							for v := 0; v < len(specialStr); v++ {
								for b := 0; b < len(specialStr); b++ {
									strs <- fmt.Sprintf("%s%s%s%s%s%s", beforeStr[i], specialStr[z], specialStr[x],
										rearStr[j], specialStr[v], specialStr[b])
									strs <- fmt.Sprintf("%s%s%s%s%s%s", strings.Title(beforeStr[i]), specialStr[z], specialStr[x],
										strings.Title(rearStr[j]), specialStr[v], specialStr[b])
								}
							}
						}
					}
				case 3:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							for c := 0; c < len(specialStr); c++ {
								for v := 0; v < len(specialStr); v++ {
									for b := 0; b < len(specialStr); b++ {
										for n := 0; n < len(specialStr); n++ {
											strs <- fmt.Sprintf("%s%s%s%s%s%s%s%s", beforeStr[i], specialStr[z], specialStr[x],
												specialStr[c], rearStr[j], specialStr[v], specialStr[b], specialStr[n])
											strs <- fmt.Sprintf("%s%s%s%s%s%s%s%s", strings.Title(beforeStr[i]), specialStr[z], specialStr[x],
												specialStr[c], strings.Title(rearStr[j]), specialStr[v], specialStr[b], specialStr[n])
										}
									}
								}
							}
						}
					}
				default:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							strs <- fmt.Sprintf("%s%s%s%s", beforeStr[i], specialStr[z], rearStr[j], specialStr[x])
							strs <- fmt.Sprintf("%s%s%s%s", strings.Title(beforeStr[i]), specialStr[z], strings.Title(rearStr[j]), specialStr[x])
						}
					}
				}
			}
		}
		wg.Done()
	}()

	// 有特殊字符 前后
	wg.Add(1)
	go func() {
		for i := 0; i < len(beforeStr); i++ {
			for j := 0; j < len(rearStr); j++ {
				switch l {
				case 2:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							for v := 0; v < len(specialStr); v++ {
								for b := 0; b < len(specialStr); b++ {
									strs <- fmt.Sprintf("%s%s%s%s%s%s", specialStr[z], specialStr[x], beforeStr[i],
										rearStr[j], specialStr[v], specialStr[b])
									strs <- fmt.Sprintf("%s%s%s%s%s%s", specialStr[z], specialStr[x], strings.Title(beforeStr[i]),
										strings.Title(rearStr[j]), specialStr[v], specialStr[b])
								}
							}
						}
					}
				case 3:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							for c := 0; c < len(specialStr); c++ {
								for v := 0; v < len(specialStr); v++ {
									for b := 0; b < len(specialStr); b++ {
										for n := 0; n < len(specialStr); n++ {
											strs <- fmt.Sprintf("%s%s%s%s%s%s%s%s", specialStr[z], specialStr[x],
												specialStr[c], beforeStr[i], rearStr[j], specialStr[v], specialStr[b], specialStr[n])
											strs <- fmt.Sprintf("%s%s%s%s%s%s%s%s", specialStr[z], specialStr[x],
												specialStr[c], strings.Title(beforeStr[i]), strings.Title(rearStr[j]), specialStr[v], specialStr[b], specialStr[n])
										}
									}
								}
							}
						}
					}
				default:
					for z := 0; z < len(specialStr); z++ {
						for x := 0; x < len(specialStr); x++ {
							strs <- fmt.Sprintf("%s%s%s%s", specialStr[z], beforeStr[i], rearStr[j], specialStr[x])
							strs <- fmt.Sprintf("%s%s%s%s", beforeStr[i], strings.Title(specialStr[z]), strings.Title(rearStr[j]), specialStr[x])
						}
					}
				}
			}
		}
		wg.Done()
	}()

	// 有特殊字符 前中后
	// wg.Add(1)
	// go func() {
	// 	for i := 0; i < len(beforeStr); i++ {
	// 		for j := 0; j < len(rearStr); j++ {
	// 			switch l {
	// 			case 2:
	// 				for z := 0; z < len(specialStr); z++ {
	// 					for x := 0; x < len(specialStr); x++ {
	// 						for v := 0; v < len(specialStr); v++ {
	// 							for b := 0; b < len(specialStr); b++ {
	// 								for n := 0; n < len(specialStr); n++ {
	// 									for m := 0; m < len(specialStr); m++ {
	// 										strs <- fmt.Sprintf("%s%s%s%s%s%s%s%s", specialStr[z], specialStr[x], beforeStr[i],
	// 											specialStr[n], specialStr[m], rearStr[j], specialStr[v], specialStr[b])
	// 									}
	// 								}
	// 							}
	// 						}
	// 					}
	// 				}
	// 			case 3:
	// 				for z := 0; z < len(specialStr); z++ {
	// 					for x := 0; x < len(specialStr); x++ {
	// 						for c := 0; c < len(specialStr); c++ {
	// 							for v := 0; v < len(specialStr); v++ {
	// 								for b := 0; b < len(specialStr); b++ {
	// 									for n := 0; n < len(specialStr); n++ {
	// 										for a := 0; a < len(specialStr); a++ {
	// 											for s := 0; s < len(specialStr); s++ {
	// 												for d := 0; d < len(specialStr); d++ {
	// 													strs <- fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s", specialStr[z], specialStr[x],
	// 														specialStr[c], beforeStr[i], specialStr[a], specialStr[s], specialStr[d], rearStr[j], specialStr[v], specialStr[b], specialStr[n])
	// 												}
	// 											}
	// 										}
	// 									}
	// 								}
	// 							}
	// 						}
	// 					}
	// 				}
	// 			default:
	// 				for z := 0; z < len(specialStr); z++ {
	// 					for x := 0; x < len(specialStr); x++ {
	// 						for c := 0; c < len(specialStr); c++ {
	// 							strs <- fmt.Sprintf("%s%s%s%s", specialStr[z], beforeStr[i], specialStr[x], rearStr[j], specialStr[c])
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// 	wg.Done()
	// }()

	wg.Wait()
}
