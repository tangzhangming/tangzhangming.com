package smsa

type phoneNumber struct{
	number int //手机号码
	iddcode int //国际区号
}


func PhoneNumber(number int, iddcode int){
	return &phoneNumber{
		number number,
		iddcode iddcode,
	}
}

// 18888888888
func (p *phoneNumber) Number(){
	return p.number
}

// +8618888888888
func (p *phoneNumber) UniversalNumber(){
	return p.number
}

// +008618888888888
func (p *phoneNumber) ZeroPrefixedNumber(){
	return p.number
}