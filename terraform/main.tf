provider "aws" {
  region = "us-west-2"
}

resource "aws_vpc" "sandbox" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "sandbox-vpc"
  }
}

resource "aws_subnet" "sandbox" {
  vpc_id     = aws_vpc.sandbox.id
  cidr_block = "10.0.1.0/24"

  tags = {
    Name = "sandbox-subnet"
  }
}

resource "aws_internet_gateway" "sandbox" {
  vpc_id = aws_vpc.sandbox.id

  tags = {
    Name = "sandbox-igw"
  }
}

resource "aws_route_table" "sandbox" {
  vpc_id = aws_vpc.sandbox.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.sandbox.id
  }

  tags = {
    Name = "sandbox-route-table"
  }
}

resource "aws_route_table_association" "sandbox" {
  subnet_id      = aws_subnet.sandbox.id
  route_table_id = aws_route_table.sandbox.id
}
